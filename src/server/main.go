package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/KubrickCode/city-bot/src/libs/discord"
	"github.com/KubrickCode/city-bot/src/server/handlers"
)

func main() {
	err := discord.Init()
	if err != nil {
		log.Fatalf("시티봇 초기화 실패: %v", err)
	}
	defer discord.Close()
	log.Println("시티봇이 실행 중입니다.")

	discord.AddHandler(handlers.TestHandler)
	discord.AddHandler(handlers.TeamCreateHandler)
	discord.AddHandler(handlers.TeamInteractionHandler)
	discord.AddHandler(handlers.TeamModalHandler)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("시티봇을 종료합니다.")
}
