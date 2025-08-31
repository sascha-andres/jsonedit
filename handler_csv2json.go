package jsonedit

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/sascha-andres/jsonedit/json/c2j"
)

// handleCSV2JSON processes a CSV file and mapping file to convert to JSON
func (app *App) handleCSV2JSON(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		app.renderCSV2JSONResult(w, r)
	}
	if r.Method == http.MethodGet {
		app.renderCSV2JSONForm(w, r)
	}
}

func (app *App) renderCSV2JSONResult(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		app.logger.Error("failed to parse form", "err", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get the CSV file from the form
	csvFile, _, err := r.FormFile("csvFile")
	if err != nil {
		app.logger.Error("failed to get CSV file from form", "err", err)
		http.Error(w, "Failed to get CSV file from form", http.StatusBadRequest)
		return
	}
	defer func() {
		err := csvFile.Close()
		if err != nil {
			app.logger.Error("failed to close CSV file", "err", err)
		}
	}()

	// Get the mapping file from the form
	mappingFile, _, err := r.FormFile("mappingFile")
	if err != nil {
		app.logger.Error("failed to get mapping file from form", "err", err)
		http.Error(w, "Failed to get mapping file from form", http.StatusBadRequest)
		return
	}
	defer func() {
		err := mappingFile.Close()
		if err != nil {
			app.logger.Error("failed to close mapping file", "err", err)
		}
	}()

	// Read the CSV file content
	csvContent, err := io.ReadAll(csvFile)
	if err != nil {
		app.logger.Error("failed to read CSV file", "err", err)
		http.Error(w, "Failed to read CSV file", http.StatusInternalServerError)
		return
	}

	// Read the mapping file content
	mappingContent, err := io.ReadAll(mappingFile)
	if err != nil {
		app.logger.Error("failed to read mapping file", "err", err)
		http.Error(w, "Failed to read mapping file", http.StatusInternalServerError)
		return
	}

	// Get form options
	options := c2j.C2JOptions{
		Array:              r.FormValue("array") == "on",
		Named:              r.FormValue("named") == "on",
		OutputType:         r.FormValue("outputType"),
		NestedPropertyName: r.FormValue("nestedPropertyName"),
		Separator:          r.FormValue("separator"),
		Logger:             app.logger.With("module", "csv2json"),
	}

	// If separator is empty, use default comma
	if options.Separator == "" {
		options.Separator = ","
	}

	// If nested property name is empty, use default "data"
	if options.NestedPropertyName == "" {
		options.NestedPropertyName = "data"
	}

	// Convert CSV to JSON
	result, contentType, err := c2j.MapCSV2JSON(options, csvContent, mappingContent)
	if err != nil {
		app.logger.Error("failed to convert CSV to JSON", "err", err)

		// Render error page
		tmpl := template.Must(template.New("csv2json").Parse(csv2jsonResultTemplate))

		data := struct {
			Result      string
			ContentType string
			Error       string
		}{
			Result:      "",
			ContentType: "",
			Error:       "Failed to convert CSV to JSON: " + err.Error(),
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			app.logger.Error("failed to render CSV2JSON result template", "err", err)
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
		}
		return
	}

	// Format the JSON for display if it's JSON
	if options.OutputType == "json" {
		var jsonObj interface{}
		err = json.Unmarshal(result, &jsonObj)
		if err != nil {
			app.logger.Error("failed to parse JSON result", "err", err)
			http.Error(w, "Failed to parse JSON result: "+err.Error(), http.StatusInternalServerError)
			return
		}

		formattedJSON, err := json.MarshalIndent(jsonObj, "", app.indent)
		if err != nil {
			app.logger.Error("failed to format JSON result", "err", err)
			http.Error(w, "Failed to format JSON result: "+err.Error(), http.StatusInternalServerError)
			return
		}
		result = formattedJSON
	}

	// Set headers for file download
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename=converted-data."+getFileExtension(options.OutputType))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(result)))

	// Write the result directly to the response
	_, err = w.Write(result)
	if err != nil {
		app.logger.Error("failed to write result to response", "err", err)
		http.Error(w, "Failed to send file", http.StatusInternalServerError)
	}
}

// getFileExtension returns the appropriate file extension based on the output type
func getFileExtension(outputType string) string {
	switch outputType {
	case "json":
		return "json"
	case "yaml":
		return "yaml"
	case "toml":
		return "toml"
	default:
		return "txt"
	}
}

func (app *App) renderCSV2JSONForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("csv2json").Parse(csv2jsonFormTemplate))
	err := tmpl.Execute(w, nil)
	if err != nil {
		app.logger.Error("failed to render CSV2JSON form template", "err", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
