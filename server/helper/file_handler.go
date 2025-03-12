package helper

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

type pdfRequest struct {
	PDFdata  []byte
	Filename string
	Password string
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

func ProcessPDFRequest(ctx fiber.Ctx) (*pdfRequest, *pdfError) {
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
		PDFdata:  pdfData,
		Filename: filename,
		Password: password,
	}, nil
}
