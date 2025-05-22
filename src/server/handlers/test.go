package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func TestHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!테스트" {
		_, err := s.ChannelMessageSend(m.ChannelID, "테스트 메시지입니다")
		if err != nil {
			log.Printf("메시지 전송 실패: %v", err)
		}
	}
}
