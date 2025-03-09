package server

import (
	"errors"
	"pdftool/types"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
	"github.com/rs/zerolog/log"
)

func errHandler(ctx fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	ua := ctx.Get(fiber.HeaderUserAgent)
	ip := ctx.IP()
	method := ctx.Method()
	path := ctx.Path()

	if ua != "" &&
		ip != "" &&
		code != fiber.StatusNotFound &&
		code != fiber.StatusMethodNotAllowed &&
		err != keyauth.ErrMissingOrMalformedAPIKey {
		log.Error().Str("UserAgent", ua).Str("IP", ip).Str("Method", method).Str("Path", path).Err(err).Send()
	}

	return ctx.Status(code).JSON(types.Response{
		Error:   true,
		Message: err.Error(),
	})
}

func authMiddleware() fiber.Handler {
	return keyauth.New(keyauth.Config{
		Next: func(ctx fiber.Ctx) bool {
			sess, err := sessionStore.Get(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get session")
				return false
			}
			defer sess.Release()

			if _, ok := sess.Get("authenticated").(bool); !ok {
				return false
			}

			return true
		},
		Validator: func(ctx fiber.Ctx, key string) (bool, error) {
			if key == types.Config.Keys.API {
				return true, nil
			}

			return false, keyauth.ErrMissingOrMalformedAPIKey
		},
		ErrorHandler: errHandler,
	})
}
