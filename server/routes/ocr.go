package routes

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"pdftool/server/helper"
	"pdftool/types"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gosimple/slug"
	"github.com/rs/zerolog/log"
)

// @Summary Perform OCR on a PDF file
// @Description Uploads a PDF file and performs OCR using Mistral OCR API
// @Tags PDF Operations
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param file formData file true "PDF file to process"
// @Success 200 {object} types.Response
// @Failure 400 {object} types.Response
// @Failure 401 {object} types.Response
// @Failure 408 {object} types.Response
// @Router /v1/ocr [post]
func OCR(ctx fiber.Ctx) error {
	done := make(chan struct{})
	defer close(done)

	file, err := ctx.FormFile("file")
	if err != nil {
		log.Error().Err(err).Caller().Send()
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusBadRequest,
			"No file uploaded",
		)
	}

	if file.Header.Get("Content-Type") != "application/pdf" {
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusBadRequest,
			"Invalid file type. Only PDF files are allowed",
		)
	}

	uploadedFile := slug.MakeLang(file.Filename, "en")
	select {
	case <-ctx.Context().Done():
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusRequestTimeout,
			"Upload cancelled",
		)
	default:
		if err := ctx.SaveFileToStorage(file, uploadedFile, types.Config.S3.Storage); err != nil {
			log.Error().Caller().Err(err).Send()
			return helper.SendErrorResponse(
				ctx,
				fiber.StatusBadRequest,
				fmt.Sprintf("failed save file to storage: %v", err),
			)
		}
	}

	mistralBody := struct {
		Model    string `json:"model"`
		Document struct {
			Type        string `json:"type"`
			DocumentURL string `json:"document_url"`
		} `json:"document"`
		IncludeImageBase64 bool `json:"include_image_base64"`
	}{
		Model: "mistral-ocr-latest",
		Document: struct {
			Type        string `json:"type"`
			DocumentURL string `json:"document_url"`
		}{
			Type:        "document_url",
			DocumentURL: fmt.Sprintf("https://%s/%s/%s", types.Config.S3.Endpoint, types.Config.S3.Bucket, uploadedFile),
		},
		IncludeImageBase64: true,
	}

	jsonBody, err := json.Marshal(mistralBody)
	if err != nil {
		log.Error().Caller().Err(err).Send()
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusBadRequest,
			fmt.Sprintf("failed to marshalling json: %v", err),
		)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", types.MistralOcrApiUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusBadRequest,
			fmt.Sprintf("failed to create request: %v", err),
		)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", types.Config.Keys.Mistral))
	req.Header.Set("User-Agent", types.AppName)

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Caller().Send()
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusBadRequest,
			fmt.Sprintf("failed to get response: %v", err),
		)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Caller().Send()
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusBadRequest,
			fmt.Sprintf("failed to read response body: %v", err),
		)
	}

	var output struct {
		Pages []struct {
			Index      int    `json:"index"`
			Markdown   string `json:"markdown"`
			Images     []any  `json:"images"`
			Dimensions struct {
				Dpi    int `json:"dpi"`
				Height int `json:"height"`
				Width  int `json:"width"`
			} `json:"dimensions"`
		} `json:"pages"`
		UsageInfo struct {
			PagesProcessed int `json:"pages_processed"`
			DocSizeBytes   int `json:"doc_size_bytes"`
		} `json:"usage_info"`
	}

	if err := json.Unmarshal(body, &output); err != nil {
		log.Error().Err(err).Caller().Send()
		return helper.SendErrorResponse(
			ctx,
			fiber.StatusBadRequest,
			fmt.Sprintf("failed to unmarshal response body: %v", err),
		)
	}

	return ctx.JSON(types.Response{
		Error: false,
		Data:  output,
	})
}
