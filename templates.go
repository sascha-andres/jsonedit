package jsonedit

// Define templates for the upload page, edit page, and comparison result page
const compareResultTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>JSON Comparison Result</title>
	<link rel="icon" type="image/png" href="/assets/favicon-96x96.png" sizes="96x96" />
	<link rel="icon" type="image/svg+xml" href="/assets/favicon.svg" />
	<link rel="shortcut icon" href="/assets/favicon.ico" />
	<link rel="apple-touch-icon" sizes="180x180" href="/assets/apple-touch-icon.png" />
	<meta name="apple-mobile-web-app-title" content="JSON edit" />
	<link rel="manifest" href="/assets/site.webmanifest" />
    <style>
        @font-face {
            font-family: 'CustomFont';
            src: url('/assets/font.ttf') format('truetype');
            font-weight: normal;
            font-style: normal;
        }
        @font-face {
            font-family: 'CustomMonoFont';
            src: url('/assets/mono.ttf') format('truetype');
            font-weight: normal;
            font-style: normal;
        }
        body { font-family: 'CustomFont', sans-serif; margin: 20px; background-color: #1E293B; color: #38BDF8; }
        h1 { color: #38BDF8; }
        pre { background-color: #0F172A; padding: 10px; border-radius: 5px; overflow: auto; color: #E2E8F0; font-family: 'CustomMonoFont', monospace; }
        button { padding: 8px 15px; background: #0F172A; color: #38BDF8; border: none; cursor: pointer; margin-right: 10px; }
        button:hover { background: #0F172A; }
        .logo-container { text-align: center; margin-bottom: 20px; }
    </style>
</head>
<body>
    <div class="logo-container">
        <img src="/assets/logo.svg" alt="JSON Edit Logo" width="200">
    </div>
    <h1>JSON Comparison Result</h1>
    <div class="comparison-result">
        {{.ComparisonResult}}
    </div>
    <div style="margin-top: 20px;">
        <button onclick="window.location.href='/'">Return to Home</button>
    </div>
</body>
</html>
`

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
    <style>
        @font-face {
            font-family: 'CustomFont';
            src: url('/assets/font.ttf') format('truetype');
            font-weight: normal;
            font-style: normal;
        }
        @font-face {
            font-family: 'CustomMonoFont';
            src: url('/assets/mono.ttf') format('truetype');
            font-weight: normal;
            font-style: normal;
        }
        body { font-family: 'CustomFont', sans-serif; margin: 20px; background-color: #1E293B; color: #38BDF8; }
        h1 { color: #38BDF8; }
        form { margin: 20px 0; }
        textarea { width: 100%; height: 300px; margin: 10px 0; font-family: 'CustomMonoFont', monospace; }
        button { padding: 8px 15px; background: #0F172A; color: #38BDF8; border: none; cursor: pointer; margin-right: 10px; }
        button:hover { background: #0F172A; }
        .error { color: red; }
        .button-container { margin: 20px 0; }
        .logo-container { text-align: center; margin-bottom: 20px; }
        input { background-color: #0F172A; color: #38BDF8; border: 1px solid #2D3748; padding: 5px; }
    </style>
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

</body>
</html>
`

const editPageTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Edit JSON</title>
	<link rel="icon" type="image/png" href="/assets/favicon-96x96.png" sizes="96x96" />
	<link rel="icon" type="image/svg+xml" href="/assets/favicon.svg" />
	<link rel="shortcut icon" href="/assets/favicon.ico" />
	<link rel="apple-touch-icon" sizes="180x180" href="/assets/apple-touch-icon.png" />
	<meta name="apple-mobile-web-app-title" content="JSON edit" />
	<link rel="manifest" href="/assets/site.webmanifest" />
    <style>
        @font-face {
            font-family: 'CustomFont';
            src: url('/assets/font.ttf') format('truetype');
            font-weight: normal;
            font-style: normal;
        }
        @font-face {
            font-family: 'CustomMonoFont';
            src: url('/assets/mono.ttf') format('truetype');
            font-weight: normal;
            font-style: normal;
        }
        body { font-family: 'CustomFont', sans-serif; margin: 20px; background-color: #1E293B; color: #38BDF8; }
        h1 { color: #38BDF8; }
        form { margin: 20px 0; }
        textarea { width: 100%; height: 300px; margin: 10px 0; font-family: 'CustomMonoFont', monospace; background-color: #0F172A; color: #38BDF8; border: 1px solid #2D3748; }
        button { padding: 8px 15px; background: #0F172A; color: #38BDF8; border: none; cursor: pointer; margin-right: 10px; }
        button:hover { background: #0F172A; }
        .delete-property-btn, .delete-array-item-btn { background: #f44336; color: white; }
        .delete-property-btn:hover, .delete-array-item-btn:hover { background: #d32f2f; color: white; }
        .error { color: red; }
        .json-field { margin: 5px 0; display: flex; align-items: center; }
        .json-field label { min-width: 150px; margin-right: 10px; font-weight: bold; }
        .json-field input { flex-grow: 1; padding: 5px; font-family: 'CustomMonoFont', monospace; background-color: #0F172A; color: #38BDF8; border: 1px solid #2D3748; }
        select { background-color: #0F172A; color: #38BDF8; border: 1px solid #2D3748; padding: 5px; }
        .hidden { display: none; }
        .logo-container { text-align: center; margin-bottom: 20px; }
    </style>
    <script>
        // Function to add a new property to an object
        function addProperty(button) {
            const path = button.getAttribute('data-path');
            const indent = parseInt(button.getAttribute('data-indent'));

            // Create a form for adding a new property
            const propertyForm = document.createElement('div');
            propertyForm.className = 'json-field';
            propertyForm.style.marginLeft = (indent * 20) + 'px';

            // Create input for property name
            const nameLabel = document.createElement('label');
            nameLabel.textContent = 'Property Name:';
            nameLabel.style.minWidth = '150px';
            nameLabel.style.marginRight = '10px';

            const nameInput = document.createElement('input');
            nameInput.type = 'text';
            nameInput.placeholder = 'Enter property name';
            nameInput.style.marginRight = '10px';

            // Create input for property value
            const valueLabel = document.createElement('label');
            valueLabel.textContent = 'Value:';
            valueLabel.style.minWidth = '50px';
            valueLabel.style.marginRight = '10px';

            const valueInput = document.createElement('input');
            valueInput.type = 'text';
            valueInput.placeholder = 'Enter value';
            valueInput.style.marginRight = '10px';

            // Create dropdown for property type
            const typeLabel = document.createElement('label');
            typeLabel.textContent = 'Type:';
            typeLabel.style.minWidth = '50px';
            typeLabel.style.marginRight = '10px';

            const typeSelect = document.createElement('select');
            typeSelect.style.marginRight = '10px';

            const typeOptions = ['string', 'number', 'boolean', 'null', 'object', 'array'];
            typeOptions.forEach(type => {
                const option = document.createElement('option');
                option.value = type;
                option.textContent = type;
                typeSelect.appendChild(option);
            });

            // Add event listener to show/hide value input based on type
            typeSelect.addEventListener('change', function() {
                const selectedType = this.value;
                if (selectedType === 'object' || selectedType === 'array') {
                    valueLabel.style.display = 'none';
                    valueInput.style.display = 'none';
                } else {
                    valueLabel.style.display = '';
                    valueInput.style.display = '';
                }
            });

            // Set default type to object and trigger change event
            typeSelect.value = 'object';
            typeSelect.dispatchEvent(new Event('change'));

            // Create add button
            const addButton = document.createElement('button');
            addButton.type = 'button';
            addButton.textContent = 'Add';
            addButton.onclick = function() {
                const propertyName = nameInput.value.trim();
                const propertyValue = valueInput.value.trim();
                const propertyType = typeSelect.value;

                if (!propertyName) {
                    alert('Property name cannot be empty');
                    return;
                }

                // Create the full path for the new property
                let fullPath = path ? path + '.' + propertyName : propertyName;

                let value;

                // Handle different types
                switch(propertyType) {
                    case 'string':
                        value = propertyValue;
                        break;
                    case 'number':
                        value = parseFloat(propertyValue) || 0;
                        break;
                    case 'boolean':
                        value = propertyValue.toLowerCase() === 'true';
                        break;
                    case 'null':
                        value = null;
                        break;
                    case 'object':
                        value = {};
                        break;
                    case 'array':
                        value = [];
                        break;
                }

                if (propertyType === 'object' || propertyType === 'array') {
                    // For objects and arrays, we'll add a nested structure
                    const newField = document.createElement('div');
                    newField.className = 'json-field';
                    newField.style.marginLeft = (indent * 20) + 'px';

                    const newLabel = document.createElement('label');
                    newLabel.setAttribute('for', fullPath);
                    newLabel.textContent = propertyName + ':';

                    newField.appendChild(newLabel);
                    newField.appendChild(document.createElement('div')); // Placeholder

                    // Insert the new field before the add property button
                    button.parentNode.parentNode.insertBefore(newField, button.parentNode);

                    // Get the JSON content and update it
                    const jsonContent = document.getElementById('jsonContent');
                    const jsonData = JSON.parse(jsonContent.value);

                    // Parse the path to access the parent object
                    const pathParts = parsePath(path);

                    // Set the new property in the JSON data
                    let parent = jsonData;
                    for (const part of pathParts) {
                        if (parent[part] === undefined) {
                            parent[part] = {};
                        }
                        parent = parent[part];
                    }
                    parent[propertyName] = value;

                    // Update the JSON content
                    jsonContent.value = JSON.stringify(jsonData, null, 4);

                    // Submit the form to update the page without triggering a POST request resend
                    const form = document.getElementById('editForm');
                    const formData = new FormData(form);

                    // Create a temporary form to submit a GET request
                    const tempForm = document.createElement('form');
                    tempForm.method = 'get';
                    tempForm.action = '/edit';

                    // Add the JSON content as a parameter
                    const input = document.createElement('input');
                    input.type = 'hidden';
                    input.name = 'jsonContent';
                    input.value = jsonContent.value;
                    tempForm.appendChild(input);

                    // Submit the form
                    document.body.appendChild(tempForm);
                    tempForm.submit();
                } else {
                    // For primitive types, add a simple input field
                    const newField = document.createElement('div');
                    newField.className = 'json-field';
                    newField.style.marginLeft = (indent * 20) + 'px';

                    const newLabel = document.createElement('label');
                    newLabel.setAttribute('for', fullPath);
                    newLabel.textContent = propertyName + ':';

                    const newInput = document.createElement('input');
                    newInput.type = 'text';
                    newInput.name = fullPath;
                    newInput.id = fullPath;
                    newInput.value = propertyValue;

                    newField.appendChild(newLabel);
                    newField.appendChild(newInput);

                    // Insert the new field before the add property button
                    button.parentNode.parentNode.insertBefore(newField, button.parentNode);
                }

                // Remove the property form
                propertyForm.remove();
            };

            // Create cancel button
            const cancelButton = document.createElement('button');
            cancelButton.type = 'button';
            cancelButton.textContent = 'Cancel';
            cancelButton.onclick = function() {
                propertyForm.remove();
            };

            // Add elements to the form
            propertyForm.appendChild(nameLabel);
            propertyForm.appendChild(nameInput);
            propertyForm.appendChild(typeLabel);
            propertyForm.appendChild(typeSelect);
            propertyForm.appendChild(valueLabel);
            propertyForm.appendChild(valueInput);
            propertyForm.appendChild(addButton);
            propertyForm.appendChild(cancelButton);

            // Insert the form before the button
            button.parentNode.parentNode.insertBefore(propertyForm, button.parentNode);
        }

        // Function to add a new item to an array
        function addArrayItem(button) {
            const path = button.getAttribute('data-path');
            const indent = parseInt(button.getAttribute('data-indent'));

            // Get the current array length
            const jsonContent = document.getElementById('jsonContent');
            const jsonData = JSON.parse(jsonContent.value);

            // Parse the path to access the array
            const pathParts = parsePath(path);

            // Get the array from the JSON data
            let array = jsonData;
            for (const part of pathParts) {
                if (array[part] === undefined) {
                    return; // Path not found
                }
                array = array[part];
            }

            // If it's not an array, convert it to one
            if (!Array.isArray(array)) {
                array = [];
                setValueByPath(jsonData, pathParts, array);
            }

            // Calculate the new index
            const newIndex = Array.isArray(array) ? array.length : 0;

            // Create a form for adding a new item
            const itemForm = document.createElement('div');
            itemForm.className = 'json-field';
            itemForm.style.marginLeft = (indent * 20) + 'px';

            // Create input for item value
            const valueLabel = document.createElement('label');
            valueLabel.textContent = 'Value:';
            valueLabel.style.minWidth = '50px';
            valueLabel.style.marginRight = '10px';

            const valueInput = document.createElement('input');
            valueInput.type = 'text';
            valueInput.placeholder = 'Enter value';
            valueInput.style.marginRight = '10px';

            // Create dropdown for item type
            const typeLabel = document.createElement('label');
            typeLabel.textContent = 'Type:';
            typeLabel.style.minWidth = '50px';
            typeLabel.style.marginRight = '10px';

            const typeSelect = document.createElement('select');
            typeSelect.style.marginRight = '10px';

            const typeOptions = ['string', 'number', 'boolean', 'null', 'object', 'array'];
            typeOptions.forEach(type => {
                const option = document.createElement('option');
                option.value = type;
                option.textContent = type;
                typeSelect.appendChild(option);
            });

            // Create add button
            const addButton = document.createElement('button');
            addButton.type = 'button';
            addButton.textContent = 'Add';
            addButton.onclick = function() {
                const itemValue = valueInput.value.trim();
                const itemType = typeSelect.value;

                // Create the full path for the new item
                let fullPath = path ? path + '[' + newIndex + ']' : '[' + newIndex + ']';

                // Create a new input field for the item
                const newField = document.createElement('div');
                newField.className = 'json-field';
                newField.style.marginLeft = (indent * 20) + 'px';

                const newLabel = document.createElement('label');
                newLabel.setAttribute('for', fullPath);
                newLabel.textContent = '[' + newIndex + ']:';

                let value;

                // Handle different types
                switch(itemType) {
                    case 'string':
                        value = itemValue;
                        break;
                    case 'number':
                        value = parseFloat(itemValue) || 0;
                        break;
                    case 'boolean':
                        value = itemValue.toLowerCase() === 'true';
                        break;
                    case 'null':
                        value = null;
                        break;
                    case 'object':
                        value = {};
                        break;
                    case 'array':
                        value = [];
                        break;
                }

                if (itemType === 'object' || itemType === 'array') {
                    // For objects and arrays, we'll add a nested structure
                    newField.appendChild(newLabel);
                    newField.appendChild(document.createElement('div')); // Placeholder

                    // Insert the new field before the add item button
                    button.parentNode.parentNode.insertBefore(newField, button.parentNode);

                    // Update the JSON data
                    setValueByPath(jsonData, [...pathParts, newIndex], value);
                    jsonContent.value = JSON.stringify(jsonData, null, 4);

                    // Submit the form to update the page without triggering a POST request resend
                    const form = document.getElementById('editForm');
                    const formData = new FormData(form);

                    // Create a temporary form to submit a GET request
                    const tempForm = document.createElement('form');
                    tempForm.method = 'get';
                    tempForm.action = '/edit';

                    // Add the JSON content as a parameter
                    const input = document.createElement('input');
                    input.type = 'hidden';
                    input.name = 'jsonContent';
                    input.value = jsonContent.value;
                    tempForm.appendChild(input);

                    // Submit the form
                    document.body.appendChild(tempForm);
                    tempForm.submit();
                } else {
                    // For primitive types, add a simple input field
                    const newInput = document.createElement('input');
                    newInput.type = 'text';
                    newInput.name = fullPath;
                    newInput.id = fullPath;
                    newInput.value = itemValue;

                    newField.appendChild(newLabel);
                    newField.appendChild(newInput);

                    // Insert the new field before the add item button
                    button.parentNode.parentNode.insertBefore(newField, button.parentNode);

                    // Update the JSON data
                    setValueByPath(jsonData, [...pathParts, newIndex], value);
                    jsonContent.value = JSON.stringify(jsonData, null, 4);
                }

                // Remove the item form
                itemForm.remove();
            };

            // Create cancel button
            const cancelButton = document.createElement('button');
            cancelButton.type = 'button';
            cancelButton.textContent = 'Cancel';
            cancelButton.onclick = function() {
                itemForm.remove();
            };

            // Add elements to the form
            itemForm.appendChild(valueLabel);
            itemForm.appendChild(valueInput);
            itemForm.appendChild(typeLabel);
            itemForm.appendChild(typeSelect);
            itemForm.appendChild(addButton);
            itemForm.appendChild(cancelButton);

            // Insert the form before the button
            button.parentNode.parentNode.insertBefore(itemForm, button.parentNode);
        }

        // Function to collect all field values and update the hidden textarea before form submission
        function updateJSONContent() {
            const form = document.getElementById('editForm');
            const inputs = form.querySelectorAll('input[type="text"]');
            const jsonContent = document.getElementById('jsonContent');

            // Start with the original JSON content
            let jsonData = JSON.parse(jsonContent.value);

            // Update each field value in the JSON data
            inputs.forEach(input => {
                const path = input.name;
                const value = input.value;

                // Skip empty paths
                if (!path) return;

                // Parse the path to access nested properties
                const pathParts = parsePath(path);

                // Update the value in the JSON data
                setValueByPath(jsonData, pathParts, parseValue(value));
            });

            // Update the hidden textarea with the new JSON content
            jsonContent.value = JSON.stringify(jsonData, null, 4);
            return true;
        }

        // Function to parse a path string into an array of path parts
        function parsePath(path) {
            const parts = [];
            let currentPart = '';
            let inBracket = false;

            for (let i = 0; i < path.length; i++) {
                const char = path[i];

                if (char === '.' && !inBracket) {
                    if (currentPart) {
                        parts.push(currentPart);
                        currentPart = '';
                    }
                } else if (char === '[') {
                    if (currentPart) {
                        parts.push(currentPart);
                        currentPart = '';
                    }
                    inBracket = true;
                } else if (char === ']') {
                    if (currentPart) {
                        parts.push(parseInt(currentPart));
                        currentPart = '';
                    }
                    inBracket = false;
                } else {
                    currentPart += char;
                }
            }

            if (currentPart) {
                parts.push(currentPart);
            }

            return parts;
        }

        // Function to set a value in a nested object by path
        function setValueByPath(obj, pathParts, value) {
            if (pathParts.length === 0) return obj;

            let current = obj;
            for (let i = 0; i < pathParts.length - 1; i++) {
                const part = pathParts[i];
                if (current[part] === undefined) {
                    current[part] = typeof pathParts[i + 1] === 'number' ? [] : {};
                }
                current = current[part];
            }

            current[pathParts[pathParts.length - 1]] = value;
            return obj;
        }

        // Function to parse a string value to the appropriate type
        function parseValue(value) {
            if (value === '') return null;
            if (value === 'true') return true;
            if (value === 'false') return false;
            if (value === 'null') return null;

            // Try to parse as number
            const num = Number(value);
            if (!isNaN(num)) return num;

            // Default to string
            return value;
        }

        // Function to delete a property from an object
        function deleteProperty(button) {
            const path = button.getAttribute('data-path');
            const key = button.getAttribute('data-key');

            // Get the JSON content
            const jsonContent = document.getElementById('jsonContent');
            const jsonData = JSON.parse(jsonContent.value);

            // Parse the path to access the parent object
            const pathParts = parsePath(path);

            // Get the parent object
            let parent = jsonData;
            for (const part of pathParts) {
                if (parent[part] === undefined) {
                    return; // Path not found
                }
                parent = parent[part];
            }

            // Delete the property
            delete parent[key];

            // Update the JSON content
            jsonContent.value = JSON.stringify(jsonData, null, 4);

            // Create a temporary form to submit a GET request
            const tempForm = document.createElement('form');
            tempForm.method = 'get';
            tempForm.action = '/edit';

            // Add the JSON content as a parameter
            const input = document.createElement('input');
            input.type = 'hidden';
            input.name = 'jsonContent';
            input.value = jsonContent.value;
            tempForm.appendChild(input);

            // Submit the form
            document.body.appendChild(tempForm);
            tempForm.submit();
        }

        // Function to delete an item from an array
        function deleteArrayItem(button) {
            const path = button.getAttribute('data-path');
            const index = parseInt(button.getAttribute('data-index'));

            // Get the JSON content
            const jsonContent = document.getElementById('jsonContent');
            const jsonData = JSON.parse(jsonContent.value);

            // Parse the path to access the parent array
            const pathParts = parsePath(path);

            // Get the parent array
            let parent = jsonData;
            for (const part of pathParts) {
                if (parent[part] === undefined) {
                    return; // Path not found
                }
                parent = parent[part];
            }

            // Check if parent is an array
            if (!Array.isArray(parent)) {
                return; // Not an array
            }

            // Remove the item at the specified index
            parent.splice(index, 1);

            // Update the JSON content
            jsonContent.value = JSON.stringify(jsonData, null, 4);

            // Create a temporary form to submit a GET request
            const tempForm = document.createElement('form');
            tempForm.method = 'get';
            tempForm.action = '/edit';

            // Add the JSON content as a parameter
            const input = document.createElement('input');
            input.type = 'hidden';
            input.name = 'jsonContent';
            input.value = jsonContent.value;
            tempForm.appendChild(input);

            // Submit the form
            document.body.appendChild(tempForm);
            tempForm.submit();
        }
    </script>
</head>
<body>
    <div class="logo-container">
        <img src="/assets/logo.svg" alt="JSON Edit Logo" width="200">
    </div>
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
</body>
</html>
`
