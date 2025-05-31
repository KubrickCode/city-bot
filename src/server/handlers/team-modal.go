package handlers

import (
	"fmt"
	"strings"

	"github.com/KubrickCode/city-bot/src/server/services"
	"github.com/bwmarrin/discordgo"
)

func TeamModalHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionModalSubmit {
		return
	}

	data := i.ModalSubmitData()

	if strings.HasPrefix(data.CustomID, "add_nickname_modal_") {
		handleNicknameModalSubmit(s, i, data)
	}
}

func handleNicknameModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ModalSubmitInteractionData) {
	messageID := strings.TrimPrefix(data.CustomID, "add_nickname_modal_")
	nicknamesInput := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	// 서비스의 비즈니스 로직 사용
	addedCount := teamService.AddMultipleParticipants(messageID, nicknamesInput, ",")

	// 원본 메시지 업데이트
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
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%d명이 추가되었습니다!", addedCount),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
