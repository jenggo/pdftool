package routes

import (
	"pdftool/server/helper"

	"github.com/gofiber/fiber/v3"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/rs/zerolog/log"
)

// @Summary Optimize a PDF file
// @Description Optimize a PDF file
// @Tags PDF Operations
// @Accept multipart/form-data,application/json
// @Produce octet-stream
// @Security ApiKeyAuth
// @Param file formData file false "PDF file to optimize"
// @Param request body object false "JSON request with base64 PDF"
// @Success 200 {file} binary
// @Failure 400 {object} types.Response
// @Failure 401 {object} types.Response
// @Failure 500 {object} types.Response
// @Router /v1/optimize [post]
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
