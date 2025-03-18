package routes

import (
	"pdftool/server/helper"

	"github.com/gofiber/fiber/v3"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/rs/zerolog/log"
)

func Optimize(ctx fiber.Ctx) error {
	result, err := helper.ProcessPDFRequest(ctx, helper.PDFProcessOptions{
		RequirePassword: false,
		OutputPrefix:    "optimized",
	})
	if err != nil {
		return helper.SendErrorResponse(ctx, err.Code, err.Message)
	}
	defer result.Cleanup()

	if err := api.ValidateFile(result.InputPath, nil); err != nil {
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusBadRequest,
			"File is invalid or corrupted. Please upload a valid PDF.",
		)
	}

	conf := model.NewDefaultConfiguration()
	conf.Optimize = true
	conf.OptimizeDuplicateContentStreams = true
	conf.OptimizeResourceDicts = true

	if err := api.OptimizeFile(result.InputPath, result.OutputPath, conf); err != nil {
		log.Error().Err(err).Caller().Send()
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusInternalServerError,
			"Failed to optimize PDF.",
		)
	}

	return ctx.Download(result.OutputPath, result.OutputName)
}
