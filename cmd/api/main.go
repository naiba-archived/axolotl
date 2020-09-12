package main

import (
	"io/ioutil"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/websocket"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/naiba/helloengineer/cmd/api/handler"
	"github.com/naiba/helloengineer/internal/model"
	"github.com/naiba/helloengineer/pkg/util"
)

var (
	oauth2config *oauth2.Config
	config       *model.Config
	cache        *bigcache.BigCache
	db           *gorm.DB
)

func init() {
	data, err := ioutil.ReadFile("data/config.yaml")
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
	db, err = gorm.Open(sqlite.Open("data/helloengineer.db"), nil)
	if err != nil {
		panic(err)
	}
	db = db.Debug()
	db.AutoMigrate(model.User{})
	util.Infof(0, "Up with config: %+v cache-cap:%d\n", config, cache.Capacity())
}

func main() {
	oauth2config = &oauth2.Config{
		ClientID:     config.GitHub.ClientID,
		ClientSecret: config.GitHub.ClientSecret,
		Endpoint:     github.Endpoint,
	}
	app := fiber.New()
	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	app.Settings.ErrorHandler = handler.DefaultError

	api := app.Group("/api")
	{
		api.Use(handler.AuthMiddleware(db))

		user := api.Group("/user")
		{
			user.Use(handler.LoginRequired(true))
			user.Get("/", handler.User)
			user.Post("/logout", handler.Logout(db))
		}

		auth := api.Group("/oauth2")
		{
			auth.Get("/login", handler.Oauth2Login(oauth2config, cache))
			auth.Get("/callback", handler.Oauth2Callback(oauth2config, cache, db))
		}
	}

	ws := app.Group("/ws")
	{
		// ws.Use(authMiddleware, requireLogin)
		ws.Get("/:meetingID", websocket.New(handler.WS()))
	}

	app.Static("/", "dist")
	app.Use(handler.NotFund)

	app.Listen(":80")
}
