package jsonedit

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/sascha-andres/jsonedit/json/fromschema"
)

// handleFromSchema processes a JSON schema file and generates an empty JSON document
func (app *App) handleFromSchema(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		app.renderJSONDocument(w, r)
	}
	if r.Method == http.MethodGet {
		app.renderJSONDocumentForm(w, r)
	}
}

// renderJSONDocumentForm renders the JSON document form on a separate page
func (app *App) renderJSONDocumentForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("compare").Parse(fromSchemaFormTemplate))
	err := tmpl.Execute(w, nil)
	if err != nil {
		app.logger.Error("failed to render upload page template", "err", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// renderJSONDocument processes a JSON schema file and generates an empty JSON document and renders the result on a separate page
func (app *App) renderJSONDocument(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		app.logger.Error("failed to parse form", "err", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("schemaFile")
	if err != nil {
		app.logger.Error("failed to get schema file from form", "err", err)
		http.Error(w, "Failed to get schema file from form", http.StatusBadRequest)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			app.logger.Error("failed to close schema file", "err", err)
		}
	}()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		app.logger.Error("failed to read schema file", "err", err)
		http.Error(w, "Failed to read schema file", http.StatusInternalServerError)
		return
	}

	// Create a schema parser
	schemaParser, err := fromschema.NewSchemaParser(app.logger.With("module", "from_schema"), content)
	if err != nil {
		app.logger.Error("failed to create schema parser", "err", err)
		http.Error(w, "Failed to parse JSON schema: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Generate an empty JSON document from the schema
	jsonData, err := schemaParser.CreateEmptyJSONDocument()
	if err != nil {
		app.logger.Error("failed to create empty JSON document", "err", err)
		http.Error(w, "Failed to create JSON document: "+err.Error(), http.StatusInternalServerError)
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
	w.Header().Set("Content-Disposition", "attachment; filename=schema_document.json")
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
