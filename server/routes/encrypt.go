package routes

import (
	"fmt"
	"os"
	"path/filepath"
	"pdftool/types"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gosimple/slug"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/rs/zerolog/log"
)

// @Summary Encrypt a PDF file
// @Description Encrypts a PDF file with password protection
// @Tags PDF Operations
// @Accept multipart/form-data
// @Produce octet-stream
// @Security ApiKeyAuth
// @Param file formData file true "PDF file to encrypt"
// @Param pdf_password formData string true "Password to encrypt the PDF"
// @Success 200 {file} binary
// @Failure 400 {object} types.Response
// @Failure 401 {object} types.Response
// @Failure 500 {object} types.Response
// @Router /v1/encrypt [post]
func Encrypt(ctx fiber.Ctx) error {
	// Get uploaded file
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Error().Err(err).Caller().Send()
		return ctx.Status(fiber.StatusBadRequest).JSON(types.Response{
			Error:   true,
			Message: "No file uploaded",
		})
	}

	if file.Header.Get("Content-Type") != "application/pdf" {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.Response{
			Error:   true,
			Message: "Invalid file type. Only PDF files are allowed",
		})
	}

	// Get password from form
	password := ctx.FormValue("pdf_password")
	if password == "" {
		log.Warn().Msgf("Password is required")
		return ctx.Status(fiber.StatusBadRequest).JSON(types.Response{
			Error:   true,
			Message: "Password is required",
		})
	}

	ext := filepath.Ext(file.Filename)
	nameWithoutExt := strings.TrimSuffix(file.Filename, ext)
	uploadedFile := slug.MakeLang(nameWithoutExt, "en") + ext

	uploadedFilePath := filepath.Join("/tmp", file.Filename)
	if err := ctx.SaveFile(file, uploadedFilePath); err != nil {
		log.Error().Err(err).Caller().Send()
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.Response{
			Error:   true,
			Message: "Failed to save uploaded file",
		})
	}
	defer os.Remove(uploadedFilePath) // Clean up temp file when done

	// Generate output filename
	outputFilename := fmt.Sprintf("encrypted_%s", uploadedFile)
	outputPath := filepath.Join("/tmp", outputFilename)

	// Configure encryption
	conf := model.NewAESConfiguration(password, password, 256)
	conf.Permissions = model.PermissionsNone

	// Encrypt the PDF
	if err := api.EncryptFile(uploadedFilePath, outputPath, conf); err != nil {
		log.Error().Err(err).Caller().Send()
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.Response{
			Error:   true,
			Message: "Failed to encrypt PDF",
		})
	}

	return ctx.Download(outputPath, outputFilename)
}
