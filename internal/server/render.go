package server

import (
	"html/template"
	"net/http"
)

type RegisterTemplateData struct {
	Username  string
	FirstName string
	Status    string
}

func (s *Server) renderTemplate(w http.ResponseWriter, filenames []string, data any) {
	t, err := template.ParseFS(templatesFs, filenames...)
	if err != nil {
		s.Logger.Error.Printf("Failed to parse templates %v: %v", filenames, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		s.Logger.Error.Printf("Failed to execute template %v: %v", filenames, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
