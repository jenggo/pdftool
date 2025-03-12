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
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/rs/zerolog/log"
)

// @Summary Encrypt a PDF file
// @Description Encrypts a PDF file with password protection
// @Tags PDF Operations
// @Accept multipart/form-data,application/json
// @Produce octet-stream
// @Security ApiKeyAuth
// @Param file formData file false "PDF file to encrypt"
// @Param request body object false "JSON request with base64 PDF"
// @Param pdf_password formData string true "Password to encrypt the PDF"
// @Success 200 {file} binary
// @Failure 400 {object} types.Response
// @Failure 401 {object} types.Response
// @Failure 500 {object} types.Response
// @Router /v1/encrypt [post]
func Encrypt(ctx fiber.Ctx) error {
	req, err := helper.ProcessPDFRequest(ctx)
	if err != nil {
		return ctx.Status(err.Code).JSON(types.Response{
			Error:   true,
			Message: err.Message,
		})
	}

	// Save PDF data to temp file
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
	outputFilename := fmt.Sprintf("encrypted_%s%s", slug.MakeLang(nameWithoutExt, "en"), ext)
	outputPath := filepath.Join("/tmp", outputFilename)
	defer os.Remove(outputPath)

	if err := api.ValidateFile(tempInput, nil); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.Response{
			Error:   true,
			Message: "File is invalid or corrupted. Please upload a valid PDF.",
		})
	}

	conf := model.NewAESConfiguration(req.Password, req.Password, 256)
	conf.Permissions = model.PermissionsNone

	if err := api.EncryptFile(tempInput, outputPath, conf); err != nil {
		log.Error().Err(err).Caller().Send()
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.Response{
			Error:   true,
			Message: "Failed to encrypt PDF",
		})
	}

	return ctx.Download(outputPath, outputFilename)
}

// @Summary Decrypt a PDF file
// @Description Decrypts a password-protected PDF file
// @Tags PDF Operations
// @Accept multipart/form-data,application/json
// @Produce octet-stream
// @Security ApiKeyAuth
// @Param file formData file false "Encrypted PDF file to decrypt"
// @Param request body object false "JSON request with base64 PDF"
// @Param pdf_password formData string true "Password to decrypt the PDF"
// @Success 200 {file} binary
// @Failure 400 {object} types.Response
// @Failure 401 {object} types.Response
// @Failure 500 {object} types.Response
// @Router /v1/decrypt [post]
func Decrypt(ctx fiber.Ctx) error {
	req, err := helper.ProcessPDFRequest(ctx)
	if err != nil {
		return ctx.Status(err.Code).JSON(types.Response{
			Error:   true,
			Message: err.Message,
		})
	}

	// Save PDF data to temp file
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
	outputFilename := fmt.Sprintf("decrypted_%s%s", slug.MakeLang(nameWithoutExt, "en"), ext)
	outputPath := filepath.Join("/tmp", outputFilename)
	defer os.Remove(outputPath)

	conf := model.NewAESConfiguration(req.Password, req.Password, 256)

	if err := api.DecryptFile(tempInput, outputPath, conf); err != nil {
		log.Error().Err(err).Caller().Send()
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.Response{
			Error:   true,
			Message: "Failed to decrypt PDF. Please check if the password is correct.",
		})
	}

	return ctx.Download(outputPath, outputFilename)
}
