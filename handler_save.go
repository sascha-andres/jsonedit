package jsonedit

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

// handleSave processes the edited JSON and provides it as a download
func (app *App) handleSave(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse form", "err", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get the JSON content from the form
	jsonContent := r.FormValue("jsonContent")

	// Validate JSON
	var jsonData interface{}
	err = json.Unmarshal([]byte(jsonContent), &jsonData)
	if err != nil {
		// Try to create a simple object with the content as a string
		simpleData := map[string]interface{}{
			"content": jsonContent,
		}
		formContent := GenerateJSONForm(app.logger, app.readOnly, simpleData, "", 0)

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
		app.logger.Error("failed to format json data", "err", err)
		http.Error(w, "Failed to format JSON", http.StatusInternalServerError)
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Disposition", "attachment; filename=document.json")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(prettyJSON)))

	// Write the JSON to the response
	_, err = w.Write(prettyJSON)
	if err != nil {
		app.logger.Error("failed to write json data", "err", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
