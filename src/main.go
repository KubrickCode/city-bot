package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("시티봇 서버가 실행 중입니다!")
	})

	err := app.Listen(":3000")
	if err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}
