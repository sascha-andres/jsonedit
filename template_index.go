package jsonedit

const indexTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>JSON Edit</title>
	<link rel="icon" type="image/png" href="/assets/favicon-96x96.png" sizes="96x96" />
	<link rel="icon" type="image/svg+xml" href="/assets/favicon.svg" />
	<link rel="shortcut icon" href="/assets/favicon.ico" />
	<link rel="apple-touch-icon" sizes="180x180" href="/assets/apple-touch-icon.png" />
	<meta name="apple-mobile-web-app-title" content="JSON edit" />
	<link rel="manifest" href="/assets/site.webmanifest" />
	<script src="/assets/htmx.js"></script>
	<script src="/assets/json-edit.js"></script>
    <link rel="stylesheet" href="assets/styles.css">
</head>
<body>
    <div class="container" id="app">
        <div class="row header" id="header">
            <!-- First row content -->
            <div class="header-left" id="header-left">
                <img src="/assets/logo.svg" alt="JSON Edit Logo" class="header-logo" width="72">
                <h1>JSON edit</h1>
            </div>
            <div class="header-right" id="header-right">
                <button id="theme-toggle" class="theme-toggle-btn">
                    <span class="light-mode-icon">☀️</span>
                    <span class="dark-mode-icon">🌙</span>
                </button>
            </div>
        </div>
        <div class="row content" id="content">
            <div class="column sidebar" id="sidebar">
                <!-- First column content -->
                <h2>Edit</h2>
				<div hx-get="/upload" hx-swap="innerHTML" hx-target="#main">Edit existing JSON</div>
				<div hx-get="/new" hx-swap="innerHTML" hx-target="#main">New JSON object</div>
				<div hx-get="/new-array" hx-swap="innerHTML" hx-target="#main">New JSON array</div>
				<h2>Functions</h2>
				<h3>Functionality</h3>
				<div hx-get="/compare" hx-swap="innerHTML" hx-target="#main">Compare documents</div>
				<div hx-get="/flatten" hx-swap="innerHTML" hx-target="#main">Flatten document</div>
				<div hx-get="/csv2json" hx-swap="innerHTML" hx-target="#main">CSV to JSON</div>
				<h3>Schema</h3>
				<div hx-get="/from-schema" hx-swap="innerHTML" hx-target="#main">JSON Document from schema</div>
				<div hx-get="/validate" hx-swap="innerHTML" hx-target="#main">Validate JSON document</div>
            </div>
            <div class="column main" id="main">
                <!-- Second column content -->
                <h2>Your selected utility will be shown here</h2>
                <p>Just select one on the left</p>
            </div>
        </div>
    </div>
</body>
</html>
`
