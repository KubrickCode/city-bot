package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/KubrickCode/city-bot/src/server/services"
	"github.com/KubrickCode/city-bot/src/server/utils"
	"github.com/bwmarrin/discordgo"
)

func TeamInteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	data := i.MessageComponentData()

	switch data.CustomID {
	case "team_join":
		handleJoinButton(s, i)
	case "team_add_nickname":
		handleAddNicknameButton(s, i)
	case "team_invite_member":
		handleInviteMemberButton(s, i)
	case "team_reset":
		handleResetButton(s, i)
	case "team_create":
		handleCreateTeamButton(s, i)
	default:
		if strings.HasPrefix(data.CustomID, "select_members_") {
			handleMemberSelect(s, i)
		} else if strings.HasPrefix(data.CustomID, "team_reshuffle_") {
			handleReshuffleButton(s, i)
		}
	}
}

func handleJoinButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	messageID := i.Message.ID
	userMention := i.Member.User.Mention()

	success := teamService.AddParticipant(messageID, userMention)
	if !success {
		utils.RespondWithError(s, i,
			"이미 참가하셨거나 인원이 가득 찼습니다!",
			"참가 실패", nil)
		return
	}

	updateTeamPanel(s, i, messageID)
}

func handleAddNicknameButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	messageID := i.Message.ID
	participants := services.TeamSessions[messageID]

	if len(participants) >= 6 {
		utils.RespondWithError(s, i,
			"이미 6명이 모두 모였습니다!",
			"닉네임 추가 시도 - 인원 초과", nil)
		return
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.TextInput{
					CustomID:    "nicknames",
					Label:       "닉네임 (여러명은 콤마로 구분)",
					Style:       discordgo.TextInputParagraph,
					Placeholder: "예: 김철수, 이영희, 박민수",
					Required:    true,
					MaxLength:   500,
				},
			},
		},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   fmt.Sprintf("add_nickname_modal_%s", messageID),
			Title:      "닉네임 추가",
			Components: components,
		},
	})

	if err != nil {
		utils.RespondWithError(s, i,
			"닉네임 입력 창을 열 수 없습니다. 다시 시도해주세요.",
			"모달 응답 실패", err)
	}
}

func handleInviteMemberButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	messageID := i.Message.ID
	participants := services.TeamSessions[messageID]

	if len(participants) >= 6 {
		utils.RespondWithError(s, i,
			"이미 6명이 모두 모였습니다!",
			"멤버 초대 시도 - 인원 초과", nil)
		return
	}

	// Guild Members를 강제로 가져오기
	members, err := s.GuildMembers(i.GuildID, "", 1000)
	if err != nil {
		utils.RespondWithError(s, i,
			"서버 멤버 정보를 가져올 수 없습니다. 잠시 후 다시 시도해주세요.",
			"길드 멤버 가져오기 실패", err)
		return
	}

	// 디버깅 로그
	log.Printf("길드 ID: %s, 가져온 멤버 수: %d", i.GuildID, len(members))

	options := createMemberSelectOptionsFromMembers(members, participants)

	// 디버깅 로그
	log.Printf("생성된 옵션 수: %d", len(options))

	if len(options) == 0 {
		utils.RespondWithError(s, i,
			fmt.Sprintf("초대할 수 있는 멤버가 없습니다. (서버 멤버: %d명)", len(members)),
			"초대 가능한 멤버 없음", nil)
		return
	}

	minValues := 1
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    fmt.Sprintf("select_members_%s", messageID),
					Placeholder: "초대할 멤버를 선택하세요",
					Options:     options,
					MinValues:   &minValues,
					MaxValues:   min(len(options), 6-len(participants)),
				},
			},
		},
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "초대할 멤버를 선택하세요:",
			Components: components,
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		utils.RespondWithError(s, i,
			"멤버 선택 메뉴를 표시할 수 없습니다. 다시 시도해주세요.",
			"멤버 선택 메뉴 응답 실패", err)
	}
}

func createMemberSelectOptionsFromMembers(members []*discordgo.Member, existingParticipants []string) []discordgo.SelectMenuOption {
	var options []discordgo.SelectMenuOption
	count := 0
	totalMembers := 0
	botMembers := 0
	existingMembers := 0

	for _, member := range members {
		totalMembers++

		if member.User.Bot {
			botMembers++
			continue
		}

		if count >= 25 {
			continue
		}

		memberMention := member.User.Mention()
		// 이미 참가한 멤버는 제외
		exists := false
		for _, participant := range existingParticipants {
			if participant == memberMention {
				exists = true
				existingMembers++
				break
			}
		}
		if exists {
			continue
		}

		displayName := member.User.Username
		if member.Nick != "" {
			displayName = member.Nick
		}

		options = append(options, discordgo.SelectMenuOption{
			Label: displayName,
			Value: memberMention,
		})
		count++
	}

	// 디버깅 로그
	log.Printf("멤버 분석 - 총: %d, 봇: %d, 이미참가: %d, 최종옵션: %d",
		totalMembers, botMembers, existingMembers, len(options))

	return options
}

func handleMemberSelect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	messageID := strings.TrimPrefix(data.CustomID, "select_members_")

	addedCount := 0
	for _, value := range data.Values {
		if teamService.AddParticipant(messageID, value) {
			addedCount++
		}
	}

	participants := services.TeamSessions[messageID]
	embed := createTeamPanelEmbed(participants)
	components := createTeamPanelComponents(len(participants))
	embeds := []*discordgo.MessageEmbed{embed}

	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    i.ChannelID,
		ID:         messageID,
		Embeds:     &embeds,
		Components: &components,
	})

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    fmt.Sprintf("%d명의 멤버가 추가되었습니다!", addedCount),
			Components: []discordgo.MessageComponent{},
		},
	})
}

func handleResetButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	messageID := i.Message.ID
	services.TeamSessions[messageID] = []string{}
	updateTeamPanel(s, i, messageID)
}

func handleCreateTeamButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	messageID := i.Message.ID
	participants := services.TeamSessions[messageID]

	if len(participants) != 6 {
		utils.RespondWithError(s, i,
			"정확히 6명이 필요합니다!",
			"팀 구성 시도 - 인원 부족", nil)
		return
	}

	team1, team2 := teamService.CreateRandomTeams(participants)
	resultEmbed := teamService.CreateTeamResultEmbed(team1, team2)

	// 다시 구성 버튼 추가
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "🔄 다시 구성",
					Style:    discordgo.PrimaryButton,
					CustomID: fmt.Sprintf("team_reshuffle_%s", strings.Join(participants, "|")),
				},
			},
		},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{resultEmbed},
			Components: components,
		},
	})

	if err != nil {
		utils.RespondWithError(s, i,
			"팀 구성 결과를 표시할 수 없습니다. 다시 시도해주세요.",
			"팀 구성 결과 응답 실패", err)
	}
}

func createMemberSelectOptions(guild *discordgo.Guild, existingParticipants []string) []discordgo.SelectMenuOption {
	var options []discordgo.SelectMenuOption
	count := 0

	for _, member := range guild.Members {
		if member.User.Bot || count >= 25 {
			continue
		}

		memberMention := member.User.Mention()
		// 이미 참가한 멤버는 제외
		exists := false
		for _, participant := range existingParticipants {
			if participant == memberMention {
				exists = true
				break
			}
		}
		if exists {
			continue
		}

		displayName := member.User.Username
		if member.Nick != "" {
			displayName = member.Nick
		}

		options = append(options, discordgo.SelectMenuOption{
			Label: displayName,
			Value: memberMention,
		})
		count++
	}

	return options
}

func handleReshuffleButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	participantsStr := strings.TrimPrefix(data.CustomID, "team_reshuffle_")
	participants := strings.Split(participantsStr, "|")

	if len(participants) != 6 {
		utils.RespondWithError(s, i,
			"참가자 정보에 오류가 있습니다.",
			"다시 구성 시도 - 참가자 정보 오류", nil)
		return
	}

	team1, team2 := teamService.CreateRandomTeams(participants)
	resultEmbed := teamService.CreateTeamResultEmbed(team1, team2)

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "🔄 다시 구성",
					Style:    discordgo.PrimaryButton,
					CustomID: fmt.Sprintf("team_reshuffle_%s", strings.Join(participants, "|")),
				},
			},
		},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{resultEmbed},
			Components: components,
		},
	})

	if err != nil {
		utils.RespondWithError(s, i,
			"팀을 다시 구성할 수 없습니다. 다시 시도해주세요.",
			"팀 다시 구성 응답 실패", err)
	}
}
