// Function to add a new property to an object
function addProperty(button) {
    const path = button.getAttribute('data-path');
    const indent = parseInt(button.getAttribute('data-indent'));

    // Create a form for adding a new property
    const propertyForm = document.createElement('div');
    propertyForm.className = 'json-field property-form';
    propertyForm.style.marginLeft = (indent * 20) + 'px';

    // Create input for property name
    const nameLabel = document.createElement('label');
    nameLabel.textContent = 'Property Name:';
    nameLabel.className = 'name-label';

    const nameInput = document.createElement('input');
    nameInput.type = 'text';
    nameInput.placeholder = 'Enter property name';
    nameInput.className = 'input-field';

    // Create input for property value
    const valueLabel = document.createElement('label');
    valueLabel.textContent = 'Value:';
    valueLabel.className = 'value-label';

    const valueInput = document.createElement('input');
    valueInput.type = 'text';
    valueInput.placeholder = 'Enter value';
    valueInput.className = 'input-field';

    // Create dropdown for property type
    const typeLabel = document.createElement('label');
    typeLabel.textContent = 'Type:';
    typeLabel.className = 'type-label';

    const typeSelect = document.createElement('select');
    typeSelect.className = 'select-field';

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
            valueLabel.classList.add('element-hidden');
            valueInput.classList.add('element-hidden');
        } else {
            valueLabel.classList.remove('element-hidden');
            valueInput.classList.remove('element-hidden');
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
    itemForm.className = 'json-field property-form';
    itemForm.style.marginLeft = (indent * 20) + 'px';

    // Create input for item value
    const valueLabel = document.createElement('label');
    valueLabel.textContent = 'Value:';
    valueLabel.className = 'value-label';

    const valueInput = document.createElement('input');
    valueInput.type = 'text';
    valueInput.placeholder = 'Enter value';
    valueInput.className = 'input-field';

    // Create dropdown for item type
    const typeLabel = document.createElement('label');
    typeLabel.textContent = 'Type:';
    typeLabel.className = 'type-label';

    const typeSelect = document.createElement('select');
    typeSelect.className = 'select-field';

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

    // Ask for confirmation before deleting
    if (!confirm('Are you sure you want to delete this property?')) {
        return; // User cancelled the operation
    }

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

    // Ask for confirmation before deleting
    if (!confirm('Are you sure you want to delete this array item?')) {
        return; // User cancelled the operation
    }

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