package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

func BuildHtml(w http.ResponseWriter, containers []DockerContainer) {
	content, err := os.ReadFile("widget.gohtml")
	if err != nil {
		slog.Error("error reading widget template", "err", err)
		return
	}

	tmpl := template.Must(template.New("webpage").Parse(string(content)))
	tmpl.Execute(w, map[string]interface{}{
		"Containers": containers,
	})
}
