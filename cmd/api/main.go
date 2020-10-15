package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/naiba/helloengineer/cmd/api/handler"
	"github.com/naiba/helloengineer/internal/model"
	"github.com/naiba/helloengineer/pkg/hub"
	"github.com/naiba/helloengineer/pkg/util"
)

const (
	PROXY_ENABLE = "PROXY_ENABLE"
)

var (
	oauth2config *oauth2.Config
	config       *model.Config
	cache        *bigcache.BigCache
	db           *gorm.DB
	frontendHost *url.URL
	pubsub       *hub.Hub
)

func init() {
	data, err := ioutil.ReadFile("data/config.yaml")
	if err != nil {
		panic(err)
	}
	frontendHost, err = url.Parse("http://localhost:8080")
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
	pubsub = hub.New()
	go pubsub.Serve()
	util.Infof(0, "Up with config: %+v cache-cap:%d\n", config, cache.Capacity())
}

func main() {
	oauth2config = &oauth2.Config{
		ClientID:     config.GitHub.ClientID,
		ClientSecret: config.GitHub.ClientSecret,
		Endpoint:     github.Endpoint,
	}
	app := fiber.New(fiber.Config{
		ErrorHandler: handler.DefaultError,
	})
	app.Use(recover.New())
	app.Use(logger.New())

	api := app.Group("/api")
	{
		api.Use(handler.AuthMiddleware(db))

		runner := api.Group("/code")
		{
			runner.Use(handler.LoginRequired(true))
			runner.Get("/list", handler.ListRunner(config))
			runner.Post("/run", handler.RunCode(config, pubsub))
		}

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
		ws.Use(handler.AuthMiddleware(db), handler.LoginRequired(true))
		ws.Get("/:meetingID", websocket.New(handler.WS(pubsub)))
	}

	if os.Getenv(PROXY_ENABLE) == "" {
		app.Static("/", "dist")
		app.Use(handler.NotFund)
	} else {
		app.Use(adaptor.HTTPHandlerFunc(func(req http.ResponseWriter, resp *http.Request) {
			httputil.NewSingleHostReverseProxy(frontendHost).ServeHTTP(req, resp)
		}))

	}

	app.Listen(":80")
}
