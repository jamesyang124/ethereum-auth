// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "James Yang",
            "email": "james_yang@htc.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/ethereum-auth/v1/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "login"
                ],
                "summary": "check signed message and authenticate user",
                "parameters": [
                    {
                        "description": "auth info for downstream auth system could carry by this field, as json format",
                        "name": "extra",
                        "in": "body",
                        "schema": {
                            "type": "object"
                        }
                    },
                    {
                        "description": "ethereum digital wallet public address",
                        "name": "paddr",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "proxy downstream authenticated response json",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/api/ethereum-auth/v1/metadata": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "metadata"
                ],
                "summary": "metadata",
                "responses": {
                    "200": {
                        "description": "metadata json ex: {\"signin-text-template\": \"Sign-in nonce: %s\", \"ttl-seconds\": 5}",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/api/ethereum-auth/v1/nonce": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/html"
                ],
                "tags": [
                    "nonce"
                ],
                "summary": "generate nonce and cached with TTL for specific chain, network, and public address",
                "parameters": [
                    {
                        "description": "ethereum digital wallet public address",
                        "name": "paddr",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "6 digit random nonce ex: 123453",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "consumes": [
                    "text/html"
                ],
                "produces": [
                    "text/html"
                ],
                "tags": [
                    "health check"
                ],
                "summary": "health check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/version": {
            "get": {
                "consumes": [
                    "text/html"
                ],
                "produces": [
                    "text/html"
                ],
                "tags": [
                    "version"
                ],
                "summary": "version",
                "responses": {
                    "200": {
                        "description": "respond current service version ex: 1.0.1",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0.0",
	Host:        "localhost:3030",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Ethereum Nonce-based Authentication Micro Service",
	Description: "Nonce-based auth with Ethereum digital wallet, datastore backed by redis",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
