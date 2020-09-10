package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gofiber/fiber"
	githubapi "github.com/google/go-github/github"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"github.com/naiba/helloengineer/internal/model"
	"github.com/naiba/helloengineer/pkg/util"
)

func User(c *fiber.Ctx) {
	c.JSON(model.Response{
		Data: c.Locals(model.KeyAuthorizedUser),
	})
}

func Logout(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) {
		user := c.Locals(model.KeyAuthorizedUser).(model.User)
		user.Sid = ""
		if err := db.Save(&user).Error; err != nil {
			c.Next(err)
			return
		}
	}
}

func Oauth2Login(config *oauth2.Config, cache *bigcache.BigCache) fiber.Handler {
	return func(c *fiber.Ctx) {
		state := util.RandStringBytesMaskImprSrcUnsafe(8)
		cache.Set(fmt.Sprintf("%s%s", model.KeyOauth2State, state), nil)
		c.Redirect(config.AuthCodeURL(state, oauth2.AccessTypeOnline), http.StatusTemporaryRedirect)
	}
}

func Oauth2Callback(config *oauth2.Config, cache *bigcache.BigCache, db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) {
		_, err := cache.Get(fmt.Sprintf("%s%s", model.KeyOauth2State, c.Query("state")))
		if err != nil {
			c.Next(err)
			return
		}
		token, err := config.Exchange(c.Context(), c.Query("code"))
		if err != nil {
			c.Next(err)
			return
		}
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token.AccessToken},
		)
		tc := oauth2.NewClient(c.Context(), ts)
		client := githubapi.NewClient(tc)
		data, _, err := client.Users.Get(c.Context(), "")
		if err != nil {
			c.Next(err)
			return
		}
		var user model.User
		if err := db.First(&user, "github_id = ?", data.GetID()).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				c.Next(err)
				return
			}
			user.GithubID = data.GetID()
			user.Nickname = data.GetLogin()
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%d%s%s%d", time.Now().UnixNano(), user.Nickname, c.IP(), user.ID)), bcrypt.MinCost)
		if err != nil {
			c.Next(err)
			return
		}
		user.Sid = string(hash)

		if err := db.Save(&user).Error; err != nil {
			c.Next(err)
			return
		}
		c.Cookie(&fiber.Cookie{
			Name:  "sid",
			Value: user.Sid,
		})

		c.Redirect("/", http.StatusTemporaryRedirect)
	}
}
