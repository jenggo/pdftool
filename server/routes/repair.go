package routes

import (
	"pdftool/server/helper"

	"github.com/gofiber/fiber/v3"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func Repair(ctx fiber.Ctx) error {
	result, err := helper.ProcessPDFRequest(ctx, helper.PDFProcessOptions{
		RequirePassword: false,
		OutputPrefix:    "repaired",
	})
	if err != nil {
		return helper.SendErrorResponse(ctx, err.Code, err.Message)
	}
	defer result.Cleanup()

	if err := api.ValidateFile(result.InputPath, nil); err == nil {
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusBadRequest,
			"File does not need to repair",
		)
	}

	if err := api.OptimizeFile(result.InputPath, result.OutputPath, nil); err != nil {
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusBadRequest,
			helper.TransformPDFCPUErrorToResponse(err),
		)
	}

	return ctx.Download(result.OutputPath, result.OutputName)
}
