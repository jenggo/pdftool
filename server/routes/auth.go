package routes

import (
	"pdftool/server/helper"
	"pdftool/types"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

var SessionStore *session.Store

func Login(ctx fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.Bind().Body(&req); err != nil {
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusBadRequest,
			"Invalid request",
		)
	}

	// Check credentials
	if req.Username != types.Config.App.Auth.User || req.Password != types.Config.App.Auth.Pass {
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusUnauthorized,
			"Invalid credentials",
		)
	}

	// Get session
	sess, err := SessionStore.Get(ctx)
	if err != nil {
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusInternalServerError,
			"Session error",
		)
	}
	defer sess.Release()

	if !sess.Fresh() {
		if err := sess.Regenerate(); err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
	}

	// Set session values
	sess.Set("authenticated", true)

	// Save session
	if err := sess.Save(); err != nil {
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusInternalServerError,
			"Could not save session",
		)
	}

	return ctx.JSON(types.Response{
		Error:   false,
		Message: "Login successful",
	})
}

func CheckAuth(ctx fiber.Ctx) error {
	sess, err := SessionStore.Get(ctx)
	if err != nil {
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusUnauthorized,
			"Unauthorized",
		)
	}
	defer sess.Release()

	authenticated := sess.Get("authenticated")
	if authenticated == nil || !authenticated.(bool) {
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusUnauthorized,
			"Unauthorized",
		)
	}

	return ctx.JSON(types.Response{
		Error:   false,
		Message: "Authenticated",
	})
}
