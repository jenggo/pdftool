package routes

import (
	"fmt"
	"os"
	"path/filepath"
	"pdftool/server/helper"
	"pdftool/types"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gosimple/slug"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/rs/zerolog/log"
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
	req, err := helper.ProcessPDFRequest(ctx)
	if err != nil {
		return ctx.Status(err.Code).JSON(types.Response{
			Error:   true,
			Message: err.Message,
		})
	}

	tempInput := filepath.Join("/tmp", req.Filename)
	if err := os.WriteFile(tempInput, req.PDFdata, 0644); err != nil {
		log.Error().Err(err).Caller().Send()
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.Response{
			Error:   true,
			Message: "Failed to process PDF",
		})
	}
	defer os.Remove(tempInput)

	ext := filepath.Ext(req.Filename)
	nameWithoutExt := strings.TrimSuffix(req.Filename, ext)
	outputFilename := fmt.Sprintf("repaired_%s%s", slug.MakeLang(nameWithoutExt, "en"), ext)
	outputPath := filepath.Join("/tmp", outputFilename)
	defer os.Remove(outputPath)

	if err := api.ValidateFile(tempInput, nil); err == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.Response{
			Error:   true,
			Message: "File does not need to repair",
		})
	}

	if err := api.OptimizeFile(tempInput, outputPath, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.Response{
			Error:   true,
			Message: fmt.Sprintf("Failed to repair file: %v", err),
		})
	}

	return ctx.Download(outputPath, outputFilename)
}
