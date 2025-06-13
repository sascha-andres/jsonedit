package jsonedit

const validateFormTemplate = `
<form action="/validate" method="post" enctype="multipart/form-data" id="form_validate" hx-encoding="multipart/form-data">
	<h2>Validate JSON against Schema</h2>
	<div>
		<label for="schemaFileValidate">JSON Schema File:</label>
		<input type="file" name="schemaFileValidate" accept=".json" required>
	</div>
	<div>
		<label for="jsonFileValidate">JSON Document:</label>
		<input type="file" name="jsonFileValidate" accept=".json" required>
	</div>
	<button form="form_validate" type="submit" hx-post="/validate" hx-swap="innerHTML" hx-target="#main">Validate</button>
</form>
`

// Define template for the validation result page
const validateResultTemplate = `
<h1>JSON Validation Result</h1>
<div class="validation-result">
	{{if .Error}}
		<div class="error-message">
			<h2>Validation Error</h2>
			<pre>{{.Error}}</pre>
		</div>
	{{else}}
		<div class="success-message">
			<h2>Valid JSON</h2>
			<p>The JSON document is valid according to the provided schema.</p>
		</div>
	{{end}}
</div>
`
