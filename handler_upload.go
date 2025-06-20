package jsonedit

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	"github.com/sascha-andres/jsonedit/json/form"
)

// handleUploadPage displays the upload form
func (app *App) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.New("upload").Parse(editPageFormTemplate))
		err := tmpl.Execute(w, nil)
		if err != nil {
			app.logger.Error("failed to render upload page template", "err", err)
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
		}
	}
	if r.Method == http.MethodPost {
		// Parse the multipart form
		err := r.ParseMultipartForm(10 << 20) // 10 MB max
		if err != nil {
			app.logger.Error("failed to parse form", "err", err)
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get the file from the form
		file, _, err := r.FormFile("jsonFile")
		if err != nil {
			app.logger.Error("failed to get file from form", "err", err)
			http.Error(w, "Failed to get file from form", http.StatusBadRequest)
			return
		}
		defer func() {
			err := file.Close()
			if err != nil {
				app.logger.Error("failed to close file", "err", err)
			}
		}()

		// Read the file content
		content, err := io.ReadAll(file)
		if err != nil {
			app.logger.Error("failed to read file from form", "err", err)
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		// Validate JSON
		var jsonData interface{}
		err = json.Unmarshal(content, &jsonData)
		if err != nil {
			// Try to create a simple object with the content as a string
			simpleData := map[string]interface{}{
				"content": string(content),
			}
			formContent := form.GenerateJSONForm(app.logger.With("module", "form"), app.readOnly, simpleData, "", 0)

			data := EditPageData{
				Content:     string(content),
				Error:       "Invalid JSON: " + err.Error(),
				FormContent: template.HTML(formContent),
				ReadOnly:    app.readOnly,
			}
			app.renderEditPage(w, data)
			return
		}

		// Pretty print the JSON
		prettyJSON, err := json.MarshalIndent(jsonData, "", app.indent)
		if err != nil {
			app.logger.Error("failed to marshal json data", "err", err)
			http.Error(w, "Failed to format JSON", http.StatusInternalServerError)
			return
		}

		// Generate form elements for each JSON field
		formContent := form.GenerateJSONForm(app.logger.With("module", "form"), app.readOnly, jsonData, "", 0)

		// Render the edit page with the JSON content and form elements
		data := EditPageData{
			Content:     string(prettyJSON),
			FormContent: template.HTML(formContent),
			ReadOnly:    app.readOnly,
		}
		app.renderEditPage(w, data)
	}
}
