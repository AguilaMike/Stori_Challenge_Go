package web

import (
	"net/http"
)

// ServeStatic configura el servidor de archivos estáticos
func ServeStatic(staticDir string) http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))
}

// GetStaticHandler devuelve un manejador para servir archivos estáticos
func GetStaticHandler(staticDir string) (string, http.Handler) {
	return "/static/", ServeStatic(staticDir)
}
