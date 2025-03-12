package routes

import (
	"pdftool/server/helper"

	"github.com/gofiber/fiber/v3"
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
	result, err := helper.ProcessPDFRequest(ctx, helper.PDFProcessOptions{
		RequirePassword: true,
		OutputPrefix:    "encrypted",
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

	conf := model.NewAESConfiguration(result.Password, result.Password, 256)
	conf.Permissions = model.PermissionsNone

	if err := api.EncryptFile(result.InputPath, result.OutputPath, conf); err != nil {
		log.Error().Err(err).Caller().Send()
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusInternalServerError,
			helper.TransformPDFCPUErrorToResponse(err),
		)
	}

	return ctx.Download(result.OutputPath, result.OutputName)
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
	result, err := helper.ProcessPDFRequest(ctx, helper.PDFProcessOptions{
		RequirePassword: true,
		OutputPrefix:    "decrypted",
	})
	if err != nil {
		return helper.SendErrorResponse(ctx, err.Code, err.Message)
	}
	defer result.Cleanup()

	conf := model.NewDefaultConfiguration()
	conf.UserPW = result.Password
	conf.OwnerPW = result.Password

	if err := api.DecryptFile(result.InputPath, result.OutputPath, conf); err != nil {
		log.Error().Err(err).Caller().Send()
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusInternalServerError,
			helper.TransformPDFCPUErrorToResponse(err),
		)
	}

	return ctx.Download(result.OutputPath, result.OutputName)
}
