package routes

import (
	"pdftool/server/helper"

	"github.com/gofiber/fiber/v3"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// @Summary Repair a PDF file
// @Description Repair a corrupt or invalid PDF file
// @Tags PDF Operations
// @Accept multipart/form-data,application/json
// @Produce octet-stream
// @Security ApiKeyAuth
// @Param file formData file false "PDF file to repair"
// @Param request body object false "JSON request with base64 PDF"
// @Success 200 {file} binary
// @Failure 400 {object} types.Response
// @Failure 401 {object} types.Response
// @Failure 500 {object} types.Response
// @Router /v1/repair [post]
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
