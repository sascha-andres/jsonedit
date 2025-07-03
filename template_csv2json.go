package jsonedit

const csv2jsonFormTemplate = `
<form action="/csv2json" method="post" enctype="multipart/form-data" id="form_csv2json" hx-encoding="multipart/form-data"/>
	<h2>Convert CSV to JSON</h2>
	<div>
		<label for="csvFile">CSV File:</label>
		<input type="file" name="csvFile" accept=".csv" required>
	</div>
	<div>
		<label for="mappingFile">Mapping File (JSON):</label>
		<input type="file" name="mappingFile" accept=".json" required>
	</div>
	<div>
		<label for="outputType">Output Format:</label>
		<select name="outputType">
			<option value="json" selected>JSON</option>
			<option value="yaml">YAML</option>
			<option value="toml">TOML</option>
		</select>
	</div>
	<div>
		<label for="separator">CSV Separator:</label>
		<input type="text" name="separator" value="," maxlength="1">
	</div>
	<div>
		<label for="nestedPropertyName">Nested Property Name (for TOML):</label>
		<input type="text" name="nestedPropertyName" value="data">
	</div>
	<div>
		<label for="array">Wrap in Array:</label>
		<input type="checkbox" name="array">
	</div>
	<div>
		<label for="named">Use Named Columns:</label>
		<input type="checkbox" name="named" checked>
	</div>
	<button form="form_csv2json" type="submit" hx-redirect="/csv2json">Convert</button>
</form>
`

// Define template for the CSV2JSON result page
const csv2jsonResultTemplate = `
<h1>CSV to JSON Conversion Result</h1>
{{if .Error}}
<div class="error-message">
	<p>Error: {{.Error}}</p>
</div>
{{else}}
<div class="result-info">
	<p>Content Type: {{.ContentType}}</p>
</div>
<div class="conversion-result">
	<pre>{{.Result}}</pre>
</div>
{{end}}
`
