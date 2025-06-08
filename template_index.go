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
	<script src="https://unpkg.com/htmx.org@2.0.4"></script>
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
				<!-- no content cxurrently -->
            </div>
        </div>
        <div class="row content" id="content">
            <div class="column sidebar" id="sidebar">
                <!-- First column content -->
                <h2>Sidebar</h2>
				<button hx-get="/new" hx-swap="innerHTML" hx-target="#main">Upload JSON</button>
                <p>This is the sidebar content with fixed width of 250px.</p>
                <!-- Adding some dummy content to demonstrate scrolling -->
                <div class="dummy-content">
                    <p>Scroll content example.</p>
                    <p>Scroll content example.</p>
                    <p>Scroll content example.</p>
                    <p>Scroll content example.</p>
                    <p>Scroll content example.</p>
                    <p>Scroll content example.</p>
                    <p>Scroll content example.</p>
                    <p>Scroll content example.</p>
                    <p>Scroll content example.</p>
                    <p>Scroll content example.</p>
                    <p>Scroll content example.</p>
                    <p>Scroll content example.</p>
                    <p>Scroll content example.</p>
                    <p>Scroll content example.</p>
                    <p>Scroll content example.</p>
                </div>
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
