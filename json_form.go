package jsonedit

import (
	"fmt"
	"html"
)

// generateJSONForm recursively generates form elements for JSON data
func (app *App) generateJSONForm(data interface{}, path string, indent int) string {
	var result string

	app.logger.Debug("If in read-only mode, render as monospaced text without edit capabilities")
	if app.readOnly {
		return app.generateReadOnlyJSON(data, path, indent)
	}

	switch v := data.(type) {
	case map[string]interface{}:
		app.logger.Debug("Handle objects")
		if len(v) == 0 {
			app.logger.Debug("Empty object - just add the button to add properties")
			objectPath := path
			result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
			result += fmt.Sprintf(`<em>Empty object</em>`)
			result += "</div>\n"

			result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
			result += fmt.Sprintf(`<button type="button" class="add-property-btn" data-path="%s" data-indent="%d" onclick="addProperty(this)">+ Add Property</button>`, objectPath, indent)
			result += "</div>\n"
		} else {
			app.logger.Debug("Handle objects with properties")
			for key, value := range v {
				fieldPath := path
				if fieldPath != "" {
					fieldPath += "." + key
				} else {
					fieldPath = key
				}

				app.logger.Debug("Add field label with indentation")
				result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
				result += fmt.Sprintf(`<label for="%s">%s:</label>`, fieldPath, key)

				app.logger.Debug("Handle nested objects and arrays differently")
				switch value.(type) {
				case map[string]interface{}, []interface{}:
					app.logger.Debug("Add delete button for nested objects and arrays")
					result += fmt.Sprintf(`<button type="button" class="delete-property-btn" data-path="%s" data-key="%s" onclick="deleteProperty(this)">Delete</button>`, path, key)
					result += "</div>\n"
					result += app.generateJSONForm(value, fieldPath, indent+1)
				default:
					app.logger.Debug("Simple value")
					strValue := ""
					if value != nil {
						strValue = fmt.Sprintf("%v", value)
					}
					result += fmt.Sprintf(`<input type="text" name="%s" id="%s" value="%s">`,
						fieldPath, fieldPath, html.EscapeString(strValue))
					app.logger.Debug("Add delete button for simple values")
					result += fmt.Sprintf(`<button type="button" class="delete-property-btn" data-path="%s" data-key="%s" onclick="deleteProperty(this)">Delete</button>`, path, key)
					result += "</div>\n"
				}
			}

			app.logger.Debug(" Add button to add new property at this level")
			objectPath := path
			result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
			result += fmt.Sprintf(`<button type="button" class="add-property-btn" data-path="%s" data-indent="%d" onclick="addProperty(this)">+ Add Property</button>`, objectPath, indent)
			result += "</div>\n"
		}
	case []interface{}:
		app.logger.Debug("Handle arrays")
		if len(v) == 0 {
			app.logger.Debug("Empty array - display a message")
			result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
			result += fmt.Sprintf(`<em>Empty array</em>`)
			result += "</div>\n"
		} else {
			app.logger.Debug("Handle arrays with items")
			for i, value := range v {
				fieldPath := fmt.Sprintf("%s[%d]", path, i)

				app.logger.Debug("Add array index label with indentation")
				result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
				result += fmt.Sprintf(`<label for="%s">[%d]:</label>`, fieldPath, i)

				app.logger.Debug("Handle nested objects and arrays differently")
				switch value.(type) {
				case map[string]interface{}, []interface{}:
					app.logger.Debug("Add delete button for nested objects and arrays in arrays")
					result += fmt.Sprintf(`<button type="button" class="delete-array-item-btn" data-path="%s" data-index="%d" onclick="deleteArrayItem(this)">Delete</button>`, path, i)
					result += "</div>\n"
					result += app.generateJSONForm(value, fieldPath, indent+1)
				default:
					app.logger.Debug("Simple value")
					strValue := ""
					if value != nil {
						strValue = fmt.Sprintf("%v", value)
					}
					result += fmt.Sprintf(`<input type="text" name="%s" id="%s" value="%s">`,
						fieldPath, fieldPath, html.EscapeString(strValue))
					app.logger.Debug("Add delete button for simple values in arrays")
					result += fmt.Sprintf(`<button type="button" class="delete-array-item-btn" data-path="%s" data-index="%d" onclick="deleteArrayItem(this)">Delete</button>`, path, i)
					result += "</div>\n"
				}
			}
		}

		app.logger.Debug("Add button to add new item to the array (for both empty and non-empty arrays)")
		arrayPath := path
		result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
		result += fmt.Sprintf(`<button type="button" class="add-array-item-btn" data-path="%s" data-indent="%d" onclick="addArrayItem(this)">+ Add Item</button>`, arrayPath, indent)
		result += "</div>\n"
	default:
		app.logger.Debug("Handle primitive values (should only happen for the root if it's not an object or array)")
		strValue := ""
		if v != nil {
			strValue = fmt.Sprintf("%v", v)
		}
		result += fmt.Sprintf(`<div class="json-field">`)
		result += fmt.Sprintf(`<label for="%s">Value:</label>`, path)
		result += fmt.Sprintf(`<input type="text" name="%s" id="%s" value="%s">`,
			path, path, html.EscapeString(strValue))
		result += "</div>\n"
	}

	return result
}

// generateReadOnlyJSON renders JSON data as monospaced text without edit capabilities
// but maintains the visual indentation and hierarchical structure
func (app *App) generateReadOnlyJSON(data interface{}, path string, indent int) string {
	var result string

	switch v := data.(type) {
	case map[string]interface{}:
		app.logger.Debug("Handle objects")
		if len(v) == 0 {
			app.logger.Debug("Empty object")
			result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
			result += fmt.Sprintf(`<em>Empty object</em>`)
			result += "</div>\n"
		} else {
			app.logger.Debug("Handle objects with properties")
			for key, value := range v {
				fieldPath := path
				if fieldPath != "" {
					fieldPath += "." + key
				} else {
					fieldPath = key
				}

				app.logger.Debug("Add field label with indentation")
				result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
				result += fmt.Sprintf(`<label>%s:</label>`, key)

				app.logger.Debug("Handle nested objects and arrays differently")
				switch value.(type) {
				case map[string]interface{}, []interface{}:
					result += "</div>\n"
					result += app.generateReadOnlyJSON(value, fieldPath, indent+1)
				default:
					app.logger.Debug("Simple value")
					strValue := ""
					if value != nil {
						strValue = fmt.Sprintf("%v", value)
					}
					result += fmt.Sprintf(`<span style="font-family: 'CustomMonoFont', monospace;">%s</span>`, html.EscapeString(strValue))
					result += "</div>\n"
				}
			}
		}
	case []interface{}:
		app.logger.Debug("Handle arrays")
		if len(v) == 0 {
			app.logger.Debug("Empty array")
			result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
			result += fmt.Sprintf(`<em>Empty array</em>`)
			result += "</div>\n"
		} else {
			app.logger.Debug("Handle arrays with items")
			for i, value := range v {
				fieldPath := fmt.Sprintf("%s[%d]", path, i)

				app.logger.Debug("Add array index label with indentation")
				result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
				result += fmt.Sprintf(`<label>[%d]:</label>`, i)

				app.logger.Debug("Handle nested objects and arrays differently")
				switch value.(type) {
				case map[string]interface{}, []interface{}:
					result += "</div>\n"
					result += app.generateReadOnlyJSON(value, fieldPath, indent+1)
				default:
					app.logger.Debug("Simple value")
					strValue := ""
					if value != nil {
						strValue = fmt.Sprintf("%v", value)
					}
					result += fmt.Sprintf(`<span style="font-family: 'CustomMonoFont', monospace;">%s</span>`, html.EscapeString(strValue))
					result += "</div>\n"
				}
			}
		}
	default:
		app.logger.Debug("Handle primitive values (should only happen for the root if it's not an object or array)")
		strValue := ""
		if v != nil {
			strValue = fmt.Sprintf("%v", v)
		}
		result += fmt.Sprintf(`<div class="json-field">`)
		result += fmt.Sprintf(`<label>Value:</label>`)
		result += fmt.Sprintf(`<span style="font-family: 'CustomMonoFont', monospace;">%s</span>`, html.EscapeString(strValue))
		result += "</div>\n"
	}

	return result
}
