package jsonedit

import (
	"encoding/json"
	"fmt"
	"html"
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
		formContent := app.generateJSONForm(simpleData, "", 0)

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
	formContent := app.generateJSONForm(jsonData, "", 0)

	// Render the edit page with the JSON content and form elements
	data := EditPageData{
		Content:     string(prettyJSON),
		FormContent: template.HTML(formContent),
		ReadOnly:    app.readOnly,
	}
	app.renderEditPage(w, data)
}

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
		formContent := app.generateJSONForm(simpleData, "", 0)

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
	formContent := app.generateJSONForm(jsonData, "", 0)

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
	formContent := app.generateJSONForm(jsonData, "", 0)

	// Render the edit page with the empty array
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

// handleFlatten processes a JSON file and flattens it
func (app *App) handleFlatten(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		app.logger.Error("failed to parse form", "err", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("jsonFileFlat")
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
		app.logger.Error("failed to read file", "err", err)
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Validate and flatten JSON
	flattenedLines, err := FlattenJSON(string(content))
	if err != nil {
		app.logger.Error("failed to flatten JSON file", "err", err)
		http.Error(w, "Failed to flatten JSON file: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare the flattened result
	var result string
	if len(flattenedLines) == 0 {
		result = "No properties found in JSON"
	} else {
		for _, line := range flattenedLines {
			result += line + "\n"
		}
	}

	// Render the flattened result on a separate page
	tmpl := template.Must(template.New("flatten").Parse(flattenResultTemplate))

	data := struct {
		FlattenResult string
	}{
		FlattenResult: result,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		app.logger.Error("failed to render flatten result template", "err", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// handleCompare processes two JSON files and compares them
func (app *App) handleCompare(w http.ResponseWriter, r *http.Request) {
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
	diff, err := app.getJSONComparison(jsonData1, jsonData2)
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
