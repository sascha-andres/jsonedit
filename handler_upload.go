package jsonedit

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
)

// handleUploadPage displays the upload form
func (app *App) handleUploadPage(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.New("upload").Parse(uploadPageTemplate))
	err := tmpl.Execute(w, nil)
	if err != nil {
		app.logger.Error("failed to render upload page template", "err", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// handleUpload processes the JSON file upload
func (app *App) handleUpload(w http.ResponseWriter, r *http.Request) {
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
		formContent := app.generateJSONForm(simpleData, "", 0)

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
	formContent := app.generateJSONForm(jsonData, "", 0)

	// Render the edit page with the JSON content and form elements
	data := EditPageData{
		Content:     string(prettyJSON),
		FormContent: template.HTML(formContent),
		ReadOnly:    app.readOnly,
	}
	app.renderEditPage(w, data)
}