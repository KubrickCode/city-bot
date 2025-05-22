package discord

import (
	"github.com/KubrickCode/city-bot/src/libs/env"
	"github.com/bwmarrin/discordgo"
)

var session *discordgo.Session

func Init() error {
	token := env.MustGetEnv("DISCORD_BOT_TOKEN")

	var err error
	session, err = discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	return session.Open()
}

func Close() {
	if session != nil {
		session.Close()
	}
}

func GetSession() *discordgo.Session {
	return session
}

func AddHandler(handler any) {
	if session != nil {
		session.AddHandler(handler)
	}
}
