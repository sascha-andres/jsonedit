package jsonedit

// Define template for the upload page
const uploadPageTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>JSON Editor</title>
	<link rel="icon" type="image/png" href="/assets/favicon-96x96.png" sizes="96x96" />
	<link rel="icon" type="image/svg+xml" href="/assets/favicon.svg" />
	<link rel="shortcut icon" href="/assets/favicon.ico" />
	<link rel="apple-touch-icon" sizes="180x180" href="/assets/apple-touch-icon.png" />
	<meta name="apple-mobile-web-app-title" content="JSON edit" />
	<link rel="manifest" href="/assets/site.webmanifest" />
    <link rel="stylesheet" href="/assets/styles.css">
</head>
<body>
    <div class="logo-container">
        <img src="/assets/logo.svg" alt="JSON Edit Logo" width="200">
    </div>
    <h1>JSON Editor</h1>
    <form action="/upload" method="post" enctype="multipart/form-data">
        <h2>Upload JSON File</h2>
        <input type="file" name="jsonFile" accept=".json" required>
        <button type="submit">Upload</button>
    </form>

    <div class="button-container">
        <h2>Or Create New JSON</h2>
        <button onclick="window.location.href='/new'">Create New Object</button>
        <button onclick="window.location.href='/new-array'">Create New Array</button>
    </div>

	<h1>JSON Compare</h1>
	<form action="/compare" method="post" enctype="multipart/form-data">
		<h2>Compare Two JSON Files</h2>
		<div>
			<label for="jsonFile1">First JSON File:</label>
			<input type="file" name="jsonFile1" accept=".json" required>
		</div>
		<div>
			<label for="jsonFile2">Second JSON File:</label>
			<input type="file" name="jsonFile2" accept=".json" required>
		</div>
		<button type="submit">Compare</button>
	</form>

	<h1>JSON Flatten</h1>
	<form action="/flatten" method="post" enctype="multipart/form-data">
		<h2>Flatten JSON File</h2>
		<div>
			<label for="jsonFileFlat">JSON File:</label>
			<input type="file" name="jsonFileFlat" accept=".json" required>
		</div>
		<button type="submit">Flatten</button>
	</form>

	<h1>Create document from schema</h1>
	<form action="/from-schema" method="post" enctype="multipart/form-data">
		<h2>Generate JSON from Schema</h2>
		<div>
			<label for="schemaFile">JSON Schema File:</label>
			<input type="file" name="schemaFile" accept=".json" required>
		</div>
		<button type="submit">Generate</button>
	</form>

</body>
</html>
`