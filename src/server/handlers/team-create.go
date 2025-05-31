package handlers

import (
	"fmt"
	"strings"

	"github.com/KubrickCode/city-bot/src/server/services"
	"github.com/KubrickCode/city-bot/src/server/utils"
	"github.com/bwmarrin/discordgo"
)

var teamService = services.NewTeamService()

func TeamCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!팀짜기" {
		createTeamPanel(s, m.ChannelID)
	}
}

func createTeamPanel(s *discordgo.Session, channelID string) {
	embed := &discordgo.MessageEmbed{
		Title:       "🏙️ 시티봇 팀 구성하기 🏙️",
		Description: "참가자 모집 중 (0/6)\n\n참가자를 선택하고 3:3 팀을 구성해보세요!",
		Color:       0x3498db,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "현재 참가자",
				Value:  "아직 없습니다",
				Inline: false,
			},
		},
	}

	components := createInitialComponents()

	msg, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})
	if err != nil {
		utils.SendChannelError(s, channelID,
			"팀 구성 패널을 생성할 수 없습니다. 잠시 후 다시 시도해주세요.",
			"팀 구성 패널 생성 실패", err)
		return
	}

	services.TeamSessions[msg.ID] = []string{}
}

func createInitialComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{Label: "참가하기", Style: discordgo.SuccessButton, CustomID: "team_join"},
				discordgo.Button{Label: "닉네임으로 추가", Style: discordgo.PrimaryButton, CustomID: "team_add_nickname"},
				discordgo.Button{Label: "서버 멤버 초대", Style: discordgo.PrimaryButton, CustomID: "team_invite_member"},
				discordgo.Button{Label: "초기화", Style: discordgo.SecondaryButton, CustomID: "team_reset"},
			},
		},
	}
}

func updateTeamPanel(s *discordgo.Session, i *discordgo.InteractionCreate, messageID string) {
	participants := services.TeamSessions[messageID]
	embed := createTeamPanelEmbed(participants)
	components := createTeamPanelComponents(len(participants))

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
}

// 팀 패널 임베드 생성
func createTeamPanelEmbed(participants []string) *discordgo.MessageEmbed {
	participantCount := len(participants)
	participantText := "아직 없습니다"

	if participantCount > 0 {
		participantText = strings.Join(participants, "\n")
	}

	return &discordgo.MessageEmbed{
		Title:       "🏙️ City-Bot 팀 구성하기 🏙️",
		Description: fmt.Sprintf("참가자 모집 중 (%d/6)\n\n참가자를 선택하고 3:3 팀을 구성해보세요!", participantCount),
		Color:       0x3498db,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "현재 참가자", Value: participantText, Inline: false},
		},
	}
}

// 팀 패널 컴포넌트 생성
func createTeamPanelComponents(participantCount int) []discordgo.MessageComponent {
	if participantCount < 6 {
		return createInitialComponents()
	} else {
		return []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{Label: "팀 구성하기!", Style: discordgo.SuccessButton, CustomID: "team_create"},
					discordgo.Button{Label: "초기화", Style: discordgo.SecondaryButton, CustomID: "team_reset"},
				},
			},
		}
	}
}

// 헬퍼 함수
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
