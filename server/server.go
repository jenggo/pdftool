package server

import (
	"pdftool/server/routes"
	"pdftool/types"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/earlydata"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/rs/zerolog/log"
)

var sessionStore *session.Store

func New() *fiber.App {
	appCfg := fiber.Config{
		AppName:           types.AppName,
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
		ErrorHandler:      errHandler,
		ProxyHeader:       "Cf-Connecting-Ip",
		StreamRequestBody: true,
	}

	if !types.Config.App.Cloudflare {
		appCfg.ProxyHeader = "X-Real-Ip"
	}

	app := fiber.New(appCfg)
	app.Use(cors.New())
	app.Use(favicon.New())
	app.Use(helmet.New())
	app.Use(earlydata.New())
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))

	sessionStore = session.NewStore(session.Config{
		CookieHTTPOnly:  true,
		CookieSecure:    true,
		CookieSameSite:  "strict",
		AbsoluteTimeout: 24 * time.Hour,
	})
	routes.SessionStore = sessionStore

	// router
	router(app)

	go func() {
		log.Log().Msgf("» %s %s listen: %s", types.AppName, types.AppVersion, types.Config.App.Listen)

		if err := app.Listen(types.Config.App.Listen, fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
			log.Error().Caller().Err(err).Send()
		}
	}()

	return app
}
