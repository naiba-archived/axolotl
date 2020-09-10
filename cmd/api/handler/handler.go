package handler

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber"
	"github.com/naiba/helloengineer/internal/bizerr"
	"github.com/naiba/helloengineer/internal/model"
	"github.com/naiba/helloengineer/pkg/util"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) {
		sid := c.Cookies("sid")
		if sid == "" || strings.TrimSpace(sid) == "" {
			c.Next()
			return
		}
		var user model.User
		if err := db.First(&user, "sid = ?", sid).Error; err != nil {
			c.Next()
			return
		}
		c.Locals(model.KeyAuthorizedUser, user)
		c.Next()
	}
}

func LoginRequired(responseJSON bool) fiber.Handler {
	return func(c *fiber.Ctx) {
		if c.Locals(model.KeyAuthorizedUser) == nil {
			if responseJSON {
				c.Next(bizerr.UnAuthorizedError)
			} else {
				c.Redirect("/", http.StatusTemporaryRedirect)
			}
			return
		}
		c.Next()
	}
}

func DefaultError(c *fiber.Ctx, err error) {
	if err, ok := err.(bizerr.BizError); ok {
		c.JSON(model.Response{
			Code: err.Code,
			Msg:  err.Error(),
		})
		return
	}
	if util.IsErrors(err, []error{
		gorm.ErrInvalidData, gorm.ErrInvalidField, gorm.ErrInvalidTransaction,
		gorm.ErrMissingWhereClause, gorm.ErrModelValueRequired, gorm.ErrModelValueRequired,
		gorm.ErrNotImplemented, gorm.ErrPrimaryKeyRequired, gorm.ErrRecordNotFound,
		gorm.ErrRegistered, gorm.ErrUnsupportedDriver, gorm.ErrUnsupportedRelation}) {
		util.Errorf(1, "gorm: %+v", err)
		c.JSON(model.Response{
			Code: bizerr.DatabaseError.Code,
			Msg:  bizerr.DatabaseError.Msg,
		})
		return
	}
	c.JSON(model.Response{
		Code: bizerr.UnknownError.Code,
		Msg:  err.Error(),
	})
	return
}

func NotFund(c *fiber.Ctx) {
	c.SendFile("dist/index.html")
}
