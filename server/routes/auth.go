package routes

import (
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
		return ctx.Status(fiber.StatusBadRequest).JSON(types.Response{
			Error:   true,
			Message: "Invalid request",
		})
	}

	// Check credentials
	if req.Username != types.Config.App.Auth.User || req.Password != types.Config.App.Auth.Pass {
		return ctx.Status(fiber.StatusUnauthorized).JSON(types.Response{
			Error:   true,
			Message: "Invalid credentials",
		})
	}

	// Get session
	sess, err := SessionStore.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.Response{
			Error:   true,
			Message: "Session error",
		})
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
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.Response{
			Error:   true,
			Message: "Could not save session",
		})
	}

	return ctx.JSON(types.Response{
		Error:   false,
		Message: "Login successful",
	})
}

func CheckAuth(ctx fiber.Ctx) error {
	sess, err := SessionStore.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(types.Response{
			Error:   true,
			Message: "Unauthorized",
		})
	}
	defer sess.Release()

	authenticated := sess.Get("authenticated")
	if authenticated == nil || !authenticated.(bool) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(types.Response{
			Error:   true,
			Message: "Unauthorized",
		})
	}

	return ctx.JSON(types.Response{
		Error:   false,
		Message: "Authenticated",
	})
}
