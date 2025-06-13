package jsonedit

import (
	"html/template"
	"io"
	"net/http"

	"github.com/sascha-andres/jsonedit/json/validate"
)

// handleValidate processes the JSON schema and document files and validates the document against the schema
func (app *App) handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		app.renderValidateResult(w, r)
	}
	if r.Method == http.MethodGet {
		app.renderValidateForm(w, r)
	}
}

// renderValidateForm renders the validation form on a separate page
func (app *App) renderValidateForm(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.New("compare").Parse(validateFormTemplate))
	err := tmpl.Execute(w, nil)
	if err != nil {
		app.logger.Error("failed to render upload page template", "err", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// renderValidateResult processes the JSON schema and document files and validates the document against the schema and renders the result on a separate page
func (app *App) renderValidateResult(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		app.logger.Error("failed to parse form", "err", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get the schema file from the form
	schemaFile, _, err := r.FormFile("schemaFileValidate")
	if err != nil {
		app.logger.Error("failed to get schema file from form", "err", err)
		http.Error(w, "Failed to get schema file from form", http.StatusBadRequest)
		return
	}
	defer func() {
		err := schemaFile.Close()
		if err != nil {
			app.logger.Error("failed to close schema file", "err", err)
		}
	}()

	// Get the JSON document file from the form
	jsonFile, _, err := r.FormFile("jsonFileValidate")
	if err != nil {
		app.logger.Error("failed to get JSON document file from form", "err", err)
		http.Error(w, "Failed to get JSON document file from form", http.StatusBadRequest)
		return
	}
	defer func() {
		err := jsonFile.Close()
		if err != nil {
			app.logger.Error("failed to close JSON document file", "err", err)
		}
	}()

	// Read the schema file content
	schemaContent, err := io.ReadAll(schemaFile)
	if err != nil {
		app.logger.Error("failed to read schema file", "err", err)
		http.Error(w, "Failed to read schema file", http.StatusInternalServerError)
		return
	}

	// Read the JSON document file content
	jsonContent, err := io.ReadAll(jsonFile)
	if err != nil {
		app.logger.Error("failed to read JSON document file", "err", err)
		http.Error(w, "Failed to read JSON document file", http.StatusInternalServerError)
		return
	}

	// Create a new JSON validator with the schema and document
	validator, err := validate.NewJSONValidator(
		validate.WithJSONSchema(schemaContent),
		validate.WithJSONDocument(jsonContent),
		validate.WithLogger(app.logger.With("module", "validate")),
	)
	if err != nil {
		app.logger.Error("failed to create JSON validator", "err", err)
		http.Error(w, "Failed to create JSON validator: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate the JSON document against the schema
	validationErr := validator.Validate()

	// Render the validation result
	tmpl := template.Must(template.New("validate").Parse(validateResultTemplate))

	data := struct {
		Error string
	}{
		Error: "",
	}

	if validationErr != nil {
		data.Error = validationErr.Error()
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		app.logger.Error("failed to render validation result template", "err", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}
