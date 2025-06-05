package jsonedit

// Define template for the flatten result page
const flattenResultTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>JSON Flatten Result</title>
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
    <h1>JSON Flatten Result</h1>
    <div class="flatten-result">
        <pre>{{.FlattenResult}}</pre>
    </div>
    <div style="margin-top: 20px;">
        <button onclick="window.location.href='/'">Return to Home</button>
    </div>
</body>
</html>
`