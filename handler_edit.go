package jsonedit

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/sascha-andres/jsonedit/json/form"
)

// handleEdit processes the JSON content from a GET request and renders the edit page
func (app *App) handleEdit(w http.ResponseWriter, r *http.Request) {
	// Get the JSON content from the query parameter
	jsonContent := r.URL.Query().Get("jsonContent")
	if jsonContent == "" {
		app.logger.Error("invalid json content")
		http.Error(w, "Missing JSON content", http.StatusBadRequest)
		return
	}

	// Validate JSON
	var jsonData interface{}
	err := json.Unmarshal([]byte(jsonContent), &jsonData)
	if err != nil {
		// Try to create a simple object with the content as a string
		simpleData := map[string]interface{}{
			"content": jsonContent,
		}
		formContent := form.GenerateJSONForm(app.logger.With("module", "form"), app.readOnly, simpleData, "", 0)

		// If we can't parse the JSON, we'll show the error and the form content
		data := EditPageData{
			Content:     jsonContent,
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

// renderEditPage renders the edit page template
func (app *App) renderEditPage(w http.ResponseWriter, data EditPageData) {
	tmpl := template.Must(template.New("edit").Parse(editPageTemplate))
	err := tmpl.Execute(w, data)
	if err != nil {
		app.logger.Error("failed to render edit page template", "err", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
