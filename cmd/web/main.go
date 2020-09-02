package main

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/template/html"
)

func main() {
	engine := html.New("assets/template", ".html")
	engine.Reload(true)
	engine.Debug(true)
	app := fiber.New(&fiber.Settings{
		Views: engine,
	})
	app.Static("/static", "assets/static")
	app.Get("/", func(c *fiber.Ctx) {
		c.Render("index", nil)
	})
	app.Listen(":80")
}
