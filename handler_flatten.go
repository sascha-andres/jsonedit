package jsonedit

import (
	"html/template"
	"io"
	"net/http"

	"github.com/sascha-andres/jsonedit/json/flatten"
)

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
	flattenedLines, err := flatten.FlattenJSON(content)
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
