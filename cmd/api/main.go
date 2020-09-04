package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	githubapi "github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gopkg.in/yaml.v3"

	"github.com/naiba/helloengineer/internal/model"
	"github.com/naiba/helloengineer/pkg/util"
)

var (
	oauth2config *oauth2.Config
	config       *model.Config
	cache        *bigcache.BigCache

	cacheKeyOauth2State = "os:"
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
	log.Printf("Up with config: %+v cache-cap:%d\n", config, cache.Capacity())
}

func main() {
	oauth2config = &oauth2.Config{
		ClientID:     config.GitHub.ClientID,
		ClientSecret: config.GitHub.ClientSecret,
		Endpoint:     github.Endpoint,
	}
	app := fiber.New()
	app.Settings.ErrorHandler = func(ctx *fiber.Ctx, err error) {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		ctx.Status(code).SendString(err.Error())
	}
	app.Use(middleware.Recover())

	api := app.Group("/api")
	{
		auth := api.Group("/oauth2")
		{
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
				user, _, err := client.Users.Get(c.Context(), "")
				if err != nil {
					c.Next(err)
					return
				}
				c.JSON(user)
			})
		}
	}

	app.Listen(":80")
}

func prefixState(state string) string {
	return fmt.Sprintf("%s%s", cacheKeyOauth2State, state)
}
