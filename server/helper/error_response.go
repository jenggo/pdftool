package helper

import (
	"pdftool/types"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func SendErrorResponse(ctx fiber.Ctx, statusCode int, message string) error {
	return ctx.Status(statusCode).JSON(types.Response{
		Error:   true,
		Message: message,
	})
}

func TransformPDFCPUErrorToResponse(err error) string {
	errMsg := strings.TrimPrefix(err.Error(), "pdfcpu: ")

	// Capitalize only the first word
	words := strings.SplitN(errMsg, " ", 2)
	if len(words) > 0 {
		firstWord := words[0]
		firstWord = strings.ToUpper(string(firstWord[0])) + firstWord[1:]
		errMsg = firstWord
		if len(words) > 1 {
			errMsg += " " + words[1]
		}
	}

	return errMsg
}
