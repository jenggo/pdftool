package helper

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gosimple/slug"
	"github.com/rs/zerolog/log"
)

type pdfRequest struct {
	InputPath  string
	OutputPath string
	OutputName string
	Password   string
	Cleanup    func()
}

type PDFProcessOptions struct {
	RequirePassword bool
	OutputPrefix    string
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

func ProcessPDFRequest(ctx fiber.Ctx, opts PDFProcessOptions) (*pdfRequest, *pdfError) {
	var (
		tempPath           string
		filename, password string
		contentType        = ctx.Get("Content-Type")
		cleanup            func()
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
		tempPath = filepath.Join("/tmp", file.Filename)
		if err := ctx.SaveFile(file, tempPath); err != nil {
			log.Error().Err(err).Caller().Send()
			return nil, newPDFError(fiber.StatusInternalServerError, "Failed to save uploaded file")
		}
		cleanup = func() {
			if err := os.Remove(tempPath); err != nil {
				log.Warn().Err(err).Msg("Failed to remove temporary file")
			}
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

		filename = request.Filename
		password = request.Password
		tempPath = filepath.Join("/tmp", "base64-"+filename)

		// Create a temp file
		tmpFile, err := os.Create(tempPath)
		if err != nil {
			log.Error().Err(err).Caller().Send()
			return nil, newPDFError(fiber.StatusInternalServerError, "Failed to create temporary file")
		}

		// Use a decoder that writes directly to file
		decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(request.Base64))
		written, err := io.Copy(tmpFile, decoder)
		tmpFile.Close() // Close immediately after writing

		if err != nil {
			os.Remove(tempPath) // Clean up on error
			log.Error().Err(err).Caller().Send()
			return nil, newPDFError(fiber.StatusBadRequest, "Invalid base64 PDF data")
		}

		if written == 0 {
			os.Remove(tempPath) // Clean up on error
			return nil, newPDFError(fiber.StatusBadRequest, "Decoded PDF data cannot be empty")
		}

		cleanup = func() {
			if err := os.Remove(tempPath); err != nil {
				log.Warn().Err(err).Msg("failed to remove temporary file")
			}
		}
	}

	// Process password requirement
	if opts.RequirePassword && password == "" {
		cleanup()
		return nil, newPDFError(fiber.StatusBadRequest, "Password is required")
	}

	// Generate output filename and path
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	outputPrefix := opts.OutputPrefix
	if outputPrefix == "" {
		outputPrefix = "processed"
	}

	outputFilename := fmt.Sprintf("%s_%s%s", outputPrefix, slug.MakeLang(nameWithoutExt, "en"), ext)
	outputPath := filepath.Join("/tmp", outputFilename)

	// Expand cleanup to handle both files
	expandedCleanup := func() {
		runtime.GC()
		cleanup()

		switch _, err := os.Stat(outputPath); {
		case err == nil:
			if err := os.Remove(outputPath); err != nil {
				log.Warn().Err(err).Msg("failed to remove output file")
			}
		case os.IsNotExist(err):
			log.Debug().Msgf("Output file %s does not exist, skipping removal", outputPath)
		default:
			log.Error().Err(err).Msgf("Error checking if output file %s exists", outputPath)
		}
	}

	return &pdfRequest{
		InputPath:  tempPath,
		OutputPath: outputPath,
		OutputName: outputFilename,
		Password:   password,
		Cleanup:    expandedCleanup,
	}, nil
}
