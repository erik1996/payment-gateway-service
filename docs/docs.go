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
        "/payment/callback/failure": {
            "get": {
                "description": "Processes a failed payment callback and redirects to a status URL.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payment"
                ],
                "summary": "Handles failed payment provider callbacks",
                "parameters": [
                    {
                        "type": "string",
                        "description": "External ID",
                        "name": "external_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "302": {
                        "description": "Redirects to status URL",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Error extracting external ID",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Failed to handle callback",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/payment/callback/success": {
            "get": {
                "description": "Processes a successful payment callback and redirects to a status URL.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payment"
                ],
                "summary": "Handles successful payment provider callbacks",
                "parameters": [
                    {
                        "type": "string",
                        "description": "External ID",
                        "name": "external_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "302": {
                        "description": "Redirects to status URL",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Error extracting external ID",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Failed to handle callback",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/payment/deposit": {
            "post": {
                "description": "Processes a deposit request and returns a URL for payment.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payment"
                ],
                "summary": "Handles deposit requests",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization token",
                        "name": "X-AUTH-TOKEN",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Validated Payment Request",
                        "name": "validatedBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/payment.PaymentRequest"
                        }
                    },
                    {
                        "description": "Example request",
                        "name": "exampleRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/payment.PaymentRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "url",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Failed to process request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/payment/withdrawal": {
            "post": {
                "description": "Processes a withdrawal request and returns a URL for payment.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payment"
                ],
                "summary": "Handles withdrawal requests",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization token",
                        "name": "X-AUTH-TOKEN",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Validated Payment Request",
                        "name": "validatedBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/payment.PaymentRequest"
                        }
                    },
                    {
                        "description": "Example request",
                        "name": "exampleRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/payment.PaymentRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "url",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Failed to process request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "payment.PaymentRequest": {
            "type": "object",
            "required": [
                "amount",
                "country_code",
                "currency_code",
                "user_id"
            ],
            "properties": {
                "amount": {
                    "type": "number"
                },
                "country_code": {
                    "type": "string"
                },
                "currency_code": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
