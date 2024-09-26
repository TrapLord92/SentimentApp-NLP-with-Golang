package docs

import (
	"bytes"
	"encoding/json"
	"fmt"
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
    "paths": {}
}`

// swaggerInfo struct holds all necessary fields for generating the swagger doc
type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it.
var SwaggerInfo = swaggerInfo{
	Version:     "1.0.0",
	Host:        "localhost",
	BasePath:    "/",
	Schemes:     []string{"http"},
	Title:       "Swagger Example API",
	Description: "This is a sample server for a Swagger API.",
}

type s struct{}

// ReadDoc generates and returns the Swagger document in JSON format.
func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.ReplaceAll(sInfo.Description, "\n", "\\n") // Clean newlines for JSON safety

	// Template to generate the swagger doc
	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, err := json.Marshal(v)
			if err != nil {
				return "[]" // Fallback to an empty array on error
			}
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		fmt.Printf("Error parsing swagger template: %v\n", err)
		return ""
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		fmt.Printf("Error executing swagger template: %v\n", err)
		return ""
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
