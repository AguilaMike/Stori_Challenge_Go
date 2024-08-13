package web

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var templates *template.Template

// InitTemplates inicializa las plantillas HTML
func InitTemplates(templateDir string) error {
	var err error
	templates, err = template.ParseGlob(filepath.Join(templateDir, "*.html"))
	return err
}

// RenderTemplate renderiza una plantilla específica con los datos proporcionados
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	return templates.ExecuteTemplate(w, tmpl, data)
}
