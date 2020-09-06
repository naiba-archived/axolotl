package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	githubapi "github.com/google/go-github/github"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/naiba/helloengineer/internal/bizerr"
	"github.com/naiba/helloengineer/internal/model"
	"github.com/naiba/helloengineer/pkg/util"
)

var (
	oauth2config *oauth2.Config
	config       *model.Config
	cache        *bigcache.BigCache
	db           *gorm.DB

	keyOauth2State    = "ko2s:"
	keyAuthorizedUser = "kau:"
)

func init() {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	config = &model.Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	cache, err = bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute * 15))
	if err != nil {
		panic(err)
	}
	db, err = gorm.Open(sqlite.Open("helloengineer.db"), nil)
	if err != nil {
		panic(err)
	}
	db = db.Debug()
	db.AutoMigrate(model.User{})
	util.Infof("Up with config: %+v cache-cap:%d\n", config, cache.Capacity())
}

func requireLogin(c *fiber.Ctx) {
	if c.Locals(keyAuthorizedUser) == nil {
		c.Next(bizerr.UnAuthorizedError)
		return
	}
	c.Next()
}

func main() {
	oauth2config = &oauth2.Config{
		ClientID:     config.GitHub.ClientID,
		ClientSecret: config.GitHub.ClientSecret,
		Endpoint:     github.Endpoint,
	}
	app := fiber.New()
	app.Use(middleware.Logger())
	app.Settings.ErrorHandler = func(c *fiber.Ctx, err error) {
		if err, ok := err.(bizerr.BizError); ok {
			c.JSON(model.Response{
				Code: err.Code,
				Msg:  err.Error(),
			})
			return
		}
		c.JSON(model.Response{
			Code: bizerr.UnknownError.Code,
			Msg:  err.Error(),
		})
		return
	}
	app.Use(middleware.Recover())

	api := app.Group("/api")
	{
		api.Use(func(c *fiber.Ctx) {
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
			c.Locals(keyAuthorizedUser, user)
			c.Next()
		})

		user := api.Group("/user")
		{
			user.Use(requireLogin)
			user.Get("/", func(c *fiber.Ctx) {
				c.JSON(model.Response{
					Data: c.Locals(keyAuthorizedUser),
				})
			})
			user.Post("/logout", func(c *fiber.Ctx) {
				user := c.Locals(keyAuthorizedUser).(model.User)
				user.Sid = ""
				if err := db.Save(&user).Error; err != nil {
					c.Next(err)
					return
				}
			})
		}

		auth := api.Group("/oauth2")
		{
			auth.Use(func(c *fiber.Ctx) {
				if c.Locals(keyAuthorizedUser) != nil {
					c.Redirect("/", http.StatusTemporaryRedirect)
					return
				}
				c.Next()
			})
			auth.Get("/login", func(c *fiber.Ctx) {
				state := util.RandStringBytesMaskImprSrcUnsafe(8)
				cache.Set(prefixState(state), nil)
				c.Redirect(oauth2config.AuthCodeURL(state, oauth2.AccessTypeOnline), http.StatusTemporaryRedirect)
			})
			auth.Get("/callback", func(c *fiber.Ctx) {
				_, err := cache.Get(prefixState(c.Query("state")))
				if err != nil {
					c.Next(err)
					return
				}
				token, err := oauth2config.Exchange(c.Context(), c.Query("code"))
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

				// TODO check if is new user

				var user model.User
				user.GithubID = data.GetID()
				user.Nickname = data.GetLogin()
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
			})
		}
	}

	app.Static("/", "dist")
	app.Use(func(c *fiber.Ctx) {
		c.SendFile("dist/index.html")
	})

	app.Listen(":80")
}

func prefixState(state string) string {
	return fmt.Sprintf("%s%s", keyOauth2State, state)
}
