package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/naiba/axolotl/internal/model"
)

func Config(conf *model.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(model.Response{
			Data: conf.Site,
		})
	}
}
