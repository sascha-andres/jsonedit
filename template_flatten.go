package jsonedit

const flattenFormTemplate = `
<form method="post" enctype="multipart/form-data" id="form_flatten" hx-encoding="multipart/form-data">
	<h2>Flatten JSON File</h2>
	<div>
		<label for="jsonFileFlat">JSON File:</label>
		<input type="file" name="jsonFileFlat" accept=".json" required>
	</div>
	<button form="form_compare" type="submit" hx-post="/flatten" hx-swap="innerHTML" hx-target="#main">Compare</button>
</form>
`

// Define template for the flatten result page
const flattenResultTemplate = `
<h1>JSON Flatten Result</h1>
<div class="flatten-result">
	<pre>{{.FlattenResult}}</pre>
</div>
<div style="margin-top: 20px;">
	<button onclick="window.location.href='/'">Return to Home</button>
</div>
`
