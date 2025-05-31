package utils

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func RespondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, userMessage string, logMessage string, err error) {
	if err != nil {
		log.Printf("%s: %v", logMessage, err)
	} else {
		log.Printf("%s", logMessage)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "❌ " + userMessage,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func SendChannelError(s *discordgo.Session, channelID, userMessage, logMessage string, err error) {
	if err != nil {
		log.Printf("%s: %v", logMessage, err)
	} else {
		log.Printf("%s", logMessage)
	}

	_, fallbackErr := s.ChannelMessageSend(channelID, "❌ "+userMessage)
	if fallbackErr != nil {
		log.Printf("에러 메시지 전송도 실패: %v", fallbackErr)
	}
}

func RespondWithSuccess(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "✅ " + message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
