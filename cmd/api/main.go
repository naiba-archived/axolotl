package main

import (
	"github.com/gofiber/fiber"
)

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) {
		c.Render("index", nil)
	})
	app.Listen(":80")
}
