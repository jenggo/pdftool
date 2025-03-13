// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/decrypt": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Decrypts a password-protected PDF file",
                "consumes": [
                    "multipart/form-data",
                    "application/json"
                ],
                "produces": [
                    "application/octet-stream"
                ],
                "tags": [
                    "PDF Operations"
                ],
                "summary": "Decrypt a PDF file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Encrypted PDF file to decrypt",
                        "name": "file",
                        "in": "formData"
                    },
                    {
                        "description": "JSON request with base64 PDF",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "type": "object"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Password to decrypt the PDF",
                        "name": "pdf_password",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        },
        "/v1/encrypt": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Encrypts a PDF file with password protection",
                "consumes": [
                    "multipart/form-data",
                    "application/json"
                ],
                "produces": [
                    "application/octet-stream"
                ],
                "tags": [
                    "PDF Operations"
                ],
                "summary": "Encrypt a PDF file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "PDF file to encrypt",
                        "name": "file",
                        "in": "formData"
                    },
                    {
                        "description": "JSON request with base64 PDF",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "type": "object"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Password to encrypt the PDF",
                        "name": "pdf_password",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        },
        "/v1/ocr": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Uploads a PDF file and performs OCR using Mistral OCR API",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "PDF Operations"
                ],
                "summary": "Perform OCR on a PDF file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "PDF file to process",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    },
                    "408": {
                        "description": "Request Timeout",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        },
        "/v1/optimize": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Optimize a PDF file",
                "consumes": [
                    "multipart/form-data",
                    "application/json"
                ],
                "produces": [
                    "application/octet-stream"
                ],
                "tags": [
                    "PDF Operations"
                ],
                "summary": "Optimize a PDF file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "PDF file to optimize",
                        "name": "file",
                        "in": "formData"
                    },
                    {
                        "description": "JSON request with base64 PDF",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        },
        "/v1/repair": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Repair a corrupt or invalid PDF file",
                "consumes": [
                    "multipart/form-data",
                    "application/json"
                ],
                "produces": [
                    "application/octet-stream"
                ],
                "tags": [
                    "PDF Operations"
                ],
                "summary": "Repair a PDF file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "PDF file to repair",
                        "name": "file",
                        "in": "formData"
                    },
                    {
                        "description": "JSON request with base64 PDF",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "types.Response": {
            "type": "object",
            "properties": {
                "data": {},
                "error": {
                    "type": "boolean"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "Enter the token with the ` + "`" + `Bearer: ` + "`" + ` prefix, e.g. \"Bearer abcde12345\".",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "localhost:2804",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "pdfTool",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
