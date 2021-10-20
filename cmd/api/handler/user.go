package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gofiber/fiber/v2"
	githubapi "github.com/google/go-github/v39/github"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"github.com/naiba/axolotl/internal/model"
	"github.com/naiba/axolotl/pkg/util"
)

func User(c *fiber.Ctx) error {
	return c.JSON(model.Response{
		Data: c.Locals(model.KeyAuthorizedUser),
	})
}

func Logout(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals(model.KeyAuthorizedUser).(model.User)
		user.Sid = ""
		if err := db.Save(&user).Error; err != nil {
			return err
		}
		return nil
	}
}

func Oauth2Login(config *oauth2.Config, cache *bigcache.BigCache) fiber.Handler {
	return func(c *fiber.Ctx) error {
		state := util.RandStringBytesMaskImprSrcUnsafe(8)
		cache.Set(fmt.Sprintf("%s%s", model.KeyOauth2State, state), nil)
		return c.Redirect(config.AuthCodeURL(state, oauth2.AccessTypeOnline), http.StatusTemporaryRedirect)
	}
}

func Oauth2Callback(config *oauth2.Config, cache *bigcache.BigCache, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := cache.Get(fmt.Sprintf("%s%s", model.KeyOauth2State, c.Query("state")))
		if err != nil {
			return err
		}
		token, err := config.Exchange(c.Context(), c.Query("code"))
		if err != nil {
			return err
		}
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token.AccessToken},
		)
		tc := oauth2.NewClient(c.Context(), ts)
		client := githubapi.NewClient(tc)
		data, _, err := client.Users.Get(c.Context(), "")
		if err != nil {
			return err
		}
		var user model.User
		if err := db.First(&user, "github_id = ?", data.GetID()).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			user.GithubID = data.GetID()
			user.Nickname = data.GetLogin()
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%d%s%s%d", time.Now().UnixNano(), user.Nickname, c.IP(), user.ID)), bcrypt.MinCost)
		if err != nil {
			return err
		}
		user.Sid = string(hash)

		if err := db.Save(&user).Error; err != nil {
			return err
		}
		c.Cookie(&fiber.Cookie{
			Name:  "sid",
			Value: user.Sid,
		})

		return c.Redirect("/", http.StatusTemporaryRedirect)
	}
}
