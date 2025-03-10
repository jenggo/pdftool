package routes

import (
	"encoding/base64"
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

type pdfRequest struct {
	pdfData  []byte
	filename string
	password string
}

type pdfError struct {
	Code    int
	Message string
}

func newPDFError(code int, message string) *pdfError {
	return &pdfError{
		Code:    code,
		Message: message,
	}
}

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
	req, err := processPDFRequest(ctx)
	if err != nil {
		return ctx.Status(err.Code).JSON(types.Response{
			Error:   true,
			Message: err.Message,
		})
	}

	// Save PDF data to temp file
	tempInput := filepath.Join("/tmp", req.filename)
	if err := os.WriteFile(tempInput, req.pdfData, 0644); err != nil {
		log.Error().Err(err).Caller().Send()
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.Response{
			Error:   true,
			Message: "Failed to process PDF",
		})
	}
	defer os.Remove(tempInput)

	ext := filepath.Ext(req.filename)
	nameWithoutExt := strings.TrimSuffix(req.filename, ext)
	outputFilename := fmt.Sprintf("encrypted_%s%s", slug.MakeLang(nameWithoutExt, "en"), ext)
	outputPath := filepath.Join("/tmp", outputFilename)

	conf := model.NewAESConfiguration(req.password, req.password, 256)
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
	req, err := processPDFRequest(ctx)
	if err != nil {
		return ctx.Status(err.Code).JSON(types.Response{
			Error:   true,
			Message: err.Message,
		})
	}

	// Save PDF data to temp file
	tempInput := filepath.Join("/tmp", req.filename)
	if err := os.WriteFile(tempInput, req.pdfData, 0644); err != nil {
		log.Error().Err(err).Caller().Send()
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.Response{
			Error:   true,
			Message: "Failed to process PDF",
		})
	}
	defer os.Remove(tempInput)

	ext := filepath.Ext(req.filename)
	nameWithoutExt := strings.TrimSuffix(req.filename, ext)
	outputFilename := fmt.Sprintf("decrypted_%s%s", slug.MakeLang(nameWithoutExt, "en"), ext)
	outputPath := filepath.Join("/tmp", outputFilename)

	conf := model.NewAESConfiguration(req.password, req.password, 256)

	if err := api.DecryptFile(tempInput, outputPath, conf); err != nil {
		log.Error().Err(err).Caller().Send()
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.Response{
			Error:   true,
			Message: "Failed to decrypt PDF. Please check if the password is correct.",
		})
	}
	defer os.Remove(outputPath)

	return ctx.Download(outputPath, outputFilename)
}

func processPDFRequest(ctx fiber.Ctx) (*pdfRequest, *pdfError) {
	var (
		pdfData            []byte
		filename, password string
		contentType        = ctx.Get("Content-Type")
	)

	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Handle multipart/form-data
		file, err := ctx.FormFile("file")
		if err != nil {
			log.Error().Err(err).Caller().Send()
			return nil, newPDFError(fiber.StatusBadRequest, "No file uploaded")
		}

		// Check if file is empty
		if file.Size == 0 {
			return nil, newPDFError(fiber.StatusBadRequest, "File cannot be empty")
		}

		if file.Header.Get("Content-Type") != "application/pdf" {
			return nil, newPDFError(fiber.StatusBadRequest, "Invalid file type. Only PDF files are allowed")
		}

		filename = file.Filename
		password = ctx.FormValue("pdf_password")

		// Save file and read its contents
		tempPath := filepath.Join("/tmp", file.Filename)
		if err := ctx.SaveFile(file, tempPath); err != nil {
			log.Error().Err(err).Caller().Send()
			return nil, newPDFError(fiber.StatusInternalServerError, "Failed to save uploaded file")
		}
		defer os.Remove(tempPath)

		pdfData, err = os.ReadFile(tempPath)
		if err != nil {
			log.Error().Err(err).Caller().Send()
			return nil, newPDFError(fiber.StatusInternalServerError, "Failed to read file")
		}

	} else {
		// Handle JSON with base64
		var request struct {
			Filename string `json:"filename"`
			Password string `json:"password"`
			Base64   string `json:"base64_pdf"`
		}

		if err := ctx.Bind().Body(&request); err != nil {
			return nil, newPDFError(fiber.StatusBadRequest, "Invalid JSON body")
		}

		// Check if base64 string is empty
		if request.Base64 == "" {
			return nil, newPDFError(fiber.StatusBadRequest, "PDF data cannot be empty")
		}

		var err error
		pdfData, err = base64.StdEncoding.DecodeString(request.Base64)
		if err != nil {
			return nil, newPDFError(fiber.StatusBadRequest, "Invalid base64 PDF data")
		}

		// Check if decoded PDF data is empty
		if len(pdfData) == 0 {
			return nil, newPDFError(fiber.StatusBadRequest, "Decoded PDF data cannot be empty")
		}

		filename = request.Filename
		password = request.Password
	}

	if password == "" {
		return nil, newPDFError(fiber.StatusBadRequest, "Password is required")
	}

	return &pdfRequest{
		pdfData:  pdfData,
		filename: filename,
		password: password,
	}, nil
}
