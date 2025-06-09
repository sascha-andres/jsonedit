package jsonedit

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io"
	"net/http"

	"github.com/sascha-andres/jsonedit/json/compare"
)

// handleCompare processes two JSON files and compares them
func (app *App) handleCompare(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		app.renderCompareResult(w, r)
	}
	if r.Method == http.MethodGet {
		app.renderCompareForm(w, r)
	}
}

func (app *App) renderCompareResult(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		app.logger.Error("failed to parse form", "err", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get the first file from the form
	file1, _, err := r.FormFile("jsonFile1")
	if err != nil {
		app.logger.Error("failed to get first file from form", "err", err)
		http.Error(w, "Failed to get first file from form", http.StatusBadRequest)
		return
	}
	defer func() {
		err := file1.Close()
		if err != nil {
			app.logger.Error("failed to close first file", "err", err)
		}
	}()

	// Get the second file from the form
	file2, _, err := r.FormFile("jsonFile2")
	if err != nil {
		app.logger.Error("failed to get second file from form", "err", err)
		http.Error(w, "Failed to get second file from form", http.StatusBadRequest)
		return
	}
	defer func() {
		err := file2.Close()
		if err != nil {
			app.logger.Error("failed to close second file", "err", err)
		}
	}()

	// Read the first file content
	content1, err := io.ReadAll(file1)
	if err != nil {
		app.logger.Error("failed to read first file", "err", err)
		http.Error(w, "Failed to read first file", http.StatusInternalServerError)
		return
	}

	// Read the second file content
	content2, err := io.ReadAll(file2)
	if err != nil {
		app.logger.Error("failed to read second file", "err", err)
		http.Error(w, "Failed to read second file", http.StatusInternalServerError)
		return
	}

	// Parse the first JSON content
	var jsonData1 interface{}
	err = json.Unmarshal(content1, &jsonData1)
	if err != nil {
		app.logger.Error("failed to parse first JSON file", "err", err)
		http.Error(w, "Failed to parse first JSON file: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the second JSON content
	var jsonData2 interface{}
	err = json.Unmarshal(content2, &jsonData2)
	if err != nil {
		app.logger.Error("failed to parse second JSON file", "err", err)
		http.Error(w, "Failed to parse second JSON file: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Compare the two JSON objects
	diff, err := compare.GetJSONComparison(jsonData1, jsonData2, app.indent)
	if err != nil {
		app.logger.Error("failed to compare JSON files", "err", err)
		http.Error(w, "Failed to compare JSON files: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the comparison result on a separate page
	tmpl := template.Must(template.New("compare").Parse(compareResultTemplate))

	var result string
	if diff == "" {
		result = "No changes"
	} else {
		result = fmt.Sprintf("<pre>%s</pre>", html.EscapeString(diff))
	}

	data := struct {
		ComparisonResult template.HTML
	}{
		ComparisonResult: template.HTML(result),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		app.logger.Error("failed to render comparison result template", "err", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
