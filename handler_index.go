package jsonedit

import (
	"html/template"
	"net/http"
)

// handleUploadPage displays the upload form
func (app *App) handleIndex(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.New("upload").Parse(indexTemplate))
	err := tmpl.Execute(w, nil)
	if err != nil {
		app.logger.Error("failed to render index template", "err", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
