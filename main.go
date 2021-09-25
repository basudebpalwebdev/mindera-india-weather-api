package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Weather API v1",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Greetings from Weather API")
	})

	app.Listen(":9999")
}
