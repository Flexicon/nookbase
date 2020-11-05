// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/search/{category}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "items"
                ],
                "summary": "Performs a search for items in the given category",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The search query",
                        "name": "q",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Item category to search through (eg: fish, insects, etc.)",
                        "name": "category",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/categories.CategoryItem"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/seasonal/{hemisphere}/{category}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "items"
                ],
                "summary": "Performs a lookup for items that are currently in season for the given hemisphere and category",
                "parameters": [
                    {
                        "enum": [
                            "northern",
                            "southern"
                        ],
                        "type": "string",
                        "description": "The hemisphere to check availability for",
                        "name": "hemisphere",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Item category to search through (eg: fish, insects, etc.)",
                        "name": "category",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/categories.CategoryItem"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "categories.CategoryItem": {
            "type": "object",
            "properties": {
                "catch_difficulty": {
                    "type": "string"
                },
                "catchphrase": {
                    "type": "string"
                },
                "critterpedia_image": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "furniture_image": {
                    "type": "string"
                },
                "icon_image": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "movement_speed": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "nh_availability": {
                    "$ref": "#/definitions/categories.YearlyAvailability"
                },
                "number": {
                    "type": "integer"
                },
                "rainy_days": {
                    "type": "boolean"
                },
                "sell": {
                    "type": "integer"
                },
                "sh_availability": {
                    "$ref": "#/definitions/categories.YearlyAvailability"
                },
                "shadow": {
                    "type": "string"
                },
                "size": {
                    "type": "string"
                },
                "spawn_rates": {
                    "type": "string"
                },
                "surface": {
                    "type": "string"
                },
                "total_catches_to_unlock": {
                    "type": "integer"
                },
                "unique_id": {
                    "type": "string"
                },
                "vision": {
                    "type": "string"
                },
                "weather": {
                    "type": "string"
                },
                "where_how": {
                    "type": "string"
                }
            }
        },
        "categories.YearlyAvailability": {
            "type": "object",
            "additionalProperties": {
                "type": "string"
            }
        },
        "main.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "invalid input"
                },
                "extra": {
                    "type": "object",
                    "additionalProperties": true
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
	Version:     "1.0",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Nookbase",
	Description: "Animal Crossing data galore",
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
	swag.Register(swag.Name, &s{})
}
