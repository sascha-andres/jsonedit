package jsonedit

import (
	"encoding/json"
	"html/template"
	"net/http"
)

// handleNewObject creates a new empty JSON object and displays it in the editor
func (app *App) handleNewObject(w http.ResponseWriter, _ *http.Request) {
	// Create an empty object
	jsonData := map[string]interface{}{}

	// Pretty print the JSON
	prettyJSON, err := json.MarshalIndent(jsonData, "", app.indent)
	if err != nil {
		app.logger.Error("failed to format json data", "err", err)
		http.Error(w, "Failed to format JSON", http.StatusInternalServerError)
		return
	}

	// Generate form elements for the empty object
	formContent := GenerateJSONForm(app.logger, app.readOnly, jsonData, "", 0)

	// Render the edit page with the empty object
	data := EditPageData{
		Content:     string(prettyJSON),
		FormContent: template.HTML(formContent),
		ReadOnly:    app.readOnly,
	}
	app.renderEditPage(w, data)
}

// handleNewArray creates a new empty JSON array and displays it in the editor
func (app *App) handleNewArray(w http.ResponseWriter, _ *http.Request) {
	// Create an empty array
	jsonData := []interface{}{}

	// Pretty print the JSON
	prettyJSON, err := json.MarshalIndent(jsonData, "", app.indent)
	if err != nil {
		app.logger.Error("failed to format json data", "err", err)
		http.Error(w, "Failed to format JSON", http.StatusInternalServerError)
		return
	}

	// Generate form elements for the empty array
	formContent := GenerateJSONForm(app.logger, app.readOnly, jsonData, "", 0)

	// Render the edit page with the empty array
	data := EditPageData{
		Content:     string(prettyJSON),
		FormContent: template.HTML(formContent),
		ReadOnly:    app.readOnly,
	}
	app.renderEditPage(w, data)
}
