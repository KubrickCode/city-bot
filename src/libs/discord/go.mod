module github.com/KubrickCode/city-bot/src/libs/discord

go 1.24.3

require (
	github.com/KubrickCode/city-bot/src/libs/env v0.0.0-00010101000000-000000000000
	github.com/bwmarrin/discordgo v0.28.1
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
)

replace github.com/KubrickCode/city-bot/src/libs/env => ../env
