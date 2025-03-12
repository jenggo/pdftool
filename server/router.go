package server

import (
	"pdftool/server/routes"
	"pdftool/types"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
	"github.com/gofiber/fiber/v3/middleware/static"
	swagger "github.com/jenggo/gofiber-swagger"
)

func router(app *fiber.App) {
	app.Get("/ping", func(ctx fiber.Ctx) error { return ctx.SendString("pong") })

	// Swagger
	if types.Config.Swagger.Enable {
		app.Get(types.Config.Swagger.Path+"/*", basicauth.New(basicauth.Config{
			Users: map[string]string{
				types.Config.App.Auth.User: types.Config.App.Auth.Pass,
			},
		}), swagger.HandlerDefault)
	}

	// UI
	app.Get("/*", static.New("ui", static.Config{Compress: true}))
	app.Get("/login", func(ctx fiber.Ctx) error { return ctx.SendFile("ui/login.html") })

	// UI Auth
	app.Post("/login", routes.Login)
	app.Get("/check-auth", routes.CheckAuth)

	// API
	v1 := app.Group("/v1", authMiddleware())
	v1.Post("/encrypt", timeoutMiddleware(2*time.Minute), routes.Encrypt)
	v1.Post("/decrypt", timeoutMiddleware(2*time.Minute), routes.Decrypt)
	v1.Post("/repair", timeoutMiddleware(2*time.Minute), routes.Repair)
	if types.Config.S3.Enable {
		v1.Post("/ocr", timeoutMiddleware(20*time.Minute), routes.OCR)
	}
}
