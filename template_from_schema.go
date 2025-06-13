package jsonedit

const fromSchemaFormTemplate = `
<form action="/from-schema" method="post" enctype="multipart/form-data" id="form_from_schema" hx-encoding="multipart/form-data">
	<h2>Generate JSON from Schema</h2>
	<div>
		<label for="schemaFile">JSON Schema File:</label>
		<input type="file" name="schemaFile" accept=".json" required>
	</div>
	<button form="form_from_schema" type="submit" hx-redirect="/from-schema">Generate</button>
</form>
`
