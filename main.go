package main

import (
	"time"

	"github.com/basudebpalwebdev/mindera-india-weather-api/api/v1/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	api := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Weather API v1",
	})
	// Enabling logging middleware
	api.Use(logger.New())

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Greetings from Weather API")
	})

	v1 := api.Group("/v1")
	// The requests hitting the route group are being cached for 3 seconds
	v1.Use(cache.New(cache.Config{
		Expiration: 3 * time.Second,
	}))

	v1.Get("/weather", handlers.FetchWeather)

	api.Listen(":9999")
}
