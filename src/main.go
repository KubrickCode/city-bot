package main

import (
	"log"

	"github.com/KubrickCode/city-bot/src/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", handlers.Index)

	err := app.Listen(":3000")
	if err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}
