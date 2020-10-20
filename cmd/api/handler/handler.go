package handler

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/naiba/axolotl/internal/bizerr"
	"github.com/naiba/axolotl/internal/model"
	"github.com/naiba/axolotl/pkg/util"
)

func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sid := c.Cookies("sid")
		if sid == "" || strings.TrimSpace(sid) == "" {
			return c.Next()
		}
		var user model.User
		if err := db.First(&user, "sid = ?", sid).Error; err != nil {
			return c.Next()
		}
		c.Locals(model.KeyAuthorizedUser, user)
		return c.Next()
	}
}

func LoginRequired(responseJSON bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals(model.KeyAuthorizedUser) == nil {
			if responseJSON {
				return bizerr.UnAuthorizedError
			}
			return c.Redirect("/", http.StatusTemporaryRedirect)
		}
		return c.Next()
	}
}

func DefaultError(c *fiber.Ctx, err error) error {
	if err, ok := err.(bizerr.BizError); ok {
		return c.JSON(model.Response{
			Code: err.Code,
			Msg:  err.Error(),
		})
	}
	if util.IsErrors(err, []error{
		gorm.ErrInvalidData, gorm.ErrInvalidField, gorm.ErrInvalidTransaction,
		gorm.ErrMissingWhereClause, gorm.ErrModelValueRequired, gorm.ErrModelValueRequired,
		gorm.ErrNotImplemented, gorm.ErrPrimaryKeyRequired, gorm.ErrRecordNotFound,
		gorm.ErrRegistered, gorm.ErrUnsupportedDriver, gorm.ErrUnsupportedRelation}) {
		util.Errorf(1, "gorm: %+v", err)
		return c.JSON(model.Response{
			Code: bizerr.DatabaseError.Code,
			Msg:  bizerr.DatabaseError.Msg,
		})
	}
	return c.JSON(model.Response{
		Code: bizerr.UnknownError.Code,
		Msg:  err.Error(),
	})
}

func NotFund(c *fiber.Ctx) error {
	return c.SendFile("dist/index.html")
}
