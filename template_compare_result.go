package jsonedit

const compareFormTemplate = `
<form action="/compare" method="post" enctype="multipart/form-data" id="form_compare" hx-encoding="multipart/form-data"/>
	<h2>Compare Two JSON Files</h2>
	<div>
		<label for="jsonFile1">First JSON File:</label>
		<input type="file" name="jsonFile1" accept=".json" required>
	</div>
	<div>
		<label for="jsonFile2">Second JSON File:</label>
		<input type="file" name="jsonFile2" accept=".json" required>
	</div>
	<button form="form_compare" type="submit" hx-post="/compare" hx-swap="innerHTML" hx-target="#main">Compare</button>
</form>
`

// Define template for the comparison result page
const compareResultTemplate = `
<h1>JSON Comparison Result</h1>
<div class="comparison-result">
	{{.ComparisonResult}}
</div>
`
