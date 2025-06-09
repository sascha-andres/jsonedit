package jsonedit

const editPageFormTemplate = `
<form action="/upload" method="post" enctype="multipart/form-data" id="form_edit" hx-encoding="multipart/form-data">
	<h2>Upload JSON File</h2>
	<input type="file" name="jsonFile" accept=".json" required>
	<button form="form_compare" type="submit" hx-post="/upload" hx-swap="innerHTML" hx-target="#main">Compare</button>
</form>
`

// Define template for the edit page
const editPageTemplate = `
<h1>Edit JSON</h1>
{{if .Error}}
<p class="error">{{.Error}}</p>
{{end}}
<form id="editForm" action="/save" method="post" onsubmit="return updateJSONContent()">
	<div id="jsonFields">
		{{.FormContent}}
	</div>
	<textarea name="jsonContent" id="jsonContent" class="hidden">{{.Content}}</textarea>
	{{if not .ReadOnly}}
	<button type="submit">Save and Download</button>
	{{end}}
	<button type="button" onclick="window.location.href='/'">Back to Upload</button>
</form>
`
