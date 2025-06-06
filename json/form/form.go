package form

import (
	"fmt"
	"html"
	"log/slog"
)

// GenerateJSONForm recursively generates form elements for JSON data
func GenerateJSONForm(logger *slog.Logger, readOnly bool, data interface{}, path string, indent int) string {
	var result string

	logger.Debug("If in read-only mode, render as monospaced text without edit capabilities")
	if readOnly {
		return GenerateReadOnlyJSON(logger, data, path, indent)
	}

	switch v := data.(type) {
	case map[string]interface{}:
		logger.Debug("Handle objects")
		if len(v) == 0 {
			logger.Debug("Empty object - just add the button to add properties")
			objectPath := path
			result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
			result += fmt.Sprintf(`<em>Empty object</em>`)
			result += "</div>\n"

			result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
			result += fmt.Sprintf(`<button type="button" class="add-property-btn" data-path="%s" data-indent="%d" onclick="addProperty(this)">+ Add Property</button>`, objectPath, indent)
			result += "</div>\n"
		} else {
			logger.Debug("Handle objects with properties")
			for key, value := range v {
				fieldPath := path
				if fieldPath != "" {
					fieldPath += "." + key
				} else {
					fieldPath = key
				}

				logger.Debug("Add field label with indentation")
				result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
				result += fmt.Sprintf(`<label for="%s">%s:</label>`, fieldPath, key)

				logger.Debug("Handle nested objects and arrays differently")
				switch value.(type) {
				case map[string]interface{}, []interface{}:
					logger.Debug("Add delete button for nested objects and arrays")
					result += fmt.Sprintf(`<button type="button" class="delete-property-btn" data-path="%s" data-key="%s" onclick="deleteProperty(this)">Delete</button>`, path, key)
					result += "</div>\n"
					result += GenerateJSONForm(logger, readOnly, value, fieldPath, indent+1)
				default:
					logger.Debug("Simple value")
					strValue := ""
					if value != nil {
						strValue = fmt.Sprintf("%v", value)
					}
					result += fmt.Sprintf(`<input type="text" name="%s" id="%s" value="%s">`,
						fieldPath, fieldPath, html.EscapeString(strValue))
					logger.Debug("Add delete button for simple values")
					result += fmt.Sprintf(`<button type="button" class="delete-property-btn" data-path="%s" data-key="%s" onclick="deleteProperty(this)">Delete</button>`, path, key)
					result += "</div>\n"
				}
			}

			logger.Debug(" Add button to add new property at this level")
			objectPath := path
			result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
			result += fmt.Sprintf(`<button type="button" class="add-property-btn" data-path="%s" data-indent="%d" onclick="addProperty(this)">+ Add Property</button>`, objectPath, indent)
			result += "</div>\n"
		}
	case []interface{}:
		logger.Debug("Handle arrays")
		if len(v) == 0 {
			logger.Debug("Empty array - display a message")
			result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
			result += fmt.Sprintf(`<em>Empty array</em>`)
			result += "</div>\n"
		} else {
			logger.Debug("Handle arrays with items")
			for i, value := range v {
				fieldPath := fmt.Sprintf("%s[%d]", path, i)

				logger.Debug("Add array index label with indentation")
				result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
				result += fmt.Sprintf(`<label for="%s">[%d]:</label>`, fieldPath, i)

				logger.Debug("Handle nested objects and arrays differently")
				switch value.(type) {
				case map[string]interface{}, []interface{}:
					logger.Debug("Add delete button for nested objects and arrays in arrays")
					result += fmt.Sprintf(`<button type="button" class="delete-array-item-btn" data-path="%s" data-index="%d" onclick="deleteArrayItem(this)">Delete</button>`, path, i)
					result += "</div>\n"
					result += GenerateJSONForm(logger, readOnly, value, fieldPath, indent+1)
				default:
					logger.Debug("Simple value")
					strValue := ""
					if value != nil {
						strValue = fmt.Sprintf("%v", value)
					}
					result += fmt.Sprintf(`<input type="text" name="%s" id="%s" value="%s">`,
						fieldPath, fieldPath, html.EscapeString(strValue))
					logger.Debug("Add delete button for simple values in arrays")
					result += fmt.Sprintf(`<button type="button" class="delete-array-item-btn" data-path="%s" data-index="%d" onclick="deleteArrayItem(this)">Delete</button>`, path, i)
					result += "</div>\n"
				}
			}
		}

		logger.Debug("Add button to add new item to the array (for both empty and non-empty arrays)")
		arrayPath := path
		result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
		result += fmt.Sprintf(`<button type="button" class="add-array-item-btn" data-path="%s" data-indent="%d" onclick="addArrayItem(this)">+ Add Item</button>`, arrayPath, indent)
		result += "</div>\n"
	default:
		logger.Debug("Handle primitive values (should only happen for the root if it's not an object or array)")
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

// GenerateReadOnlyJSON renders JSON data as monospaced text without edit capabilities
// but maintains the visual indentation and hierarchical structure
func GenerateReadOnlyJSON(logger *slog.Logger, data interface{}, path string, indent int) string {
	var result string

	switch v := data.(type) {
	case map[string]interface{}:
		logger.Debug("Handle objects")
		if len(v) == 0 {
			logger.Debug("Empty object")
			result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
			result += fmt.Sprintf(`<em>Empty object</em>`)
			result += "</div>\n"
		} else {
			logger.Debug("Handle objects with properties")
			for key, value := range v {
				fieldPath := path
				if fieldPath != "" {
					fieldPath += "." + key
				} else {
					fieldPath = key
				}

				logger.Debug("Add field label with indentation")
				result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
				result += fmt.Sprintf(`<label>%s:</label>`, key)

				logger.Debug("Handle nested objects and arrays differently")
				switch value.(type) {
				case map[string]interface{}, []interface{}:
					result += "</div>\n"
					result += GenerateReadOnlyJSON(logger, value, fieldPath, indent+1)
				default:
					logger.Debug("Simple value")
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
		logger.Debug("Handle arrays")
		if len(v) == 0 {
			logger.Debug("Empty array")
			result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
			result += fmt.Sprintf(`<em>Empty array</em>`)
			result += "</div>\n"
		} else {
			logger.Debug("Handle arrays with items")
			for i, value := range v {
				fieldPath := fmt.Sprintf("%s[%d]", path, i)

				logger.Debug("Add array index label with indentation")
				result += fmt.Sprintf(`<div class="json-field" style="margin-left: %dpx;">`, indent*20)
				result += fmt.Sprintf(`<label>[%d]:</label>`, i)

				logger.Debug("Handle nested objects and arrays differently")
				switch value.(type) {
				case map[string]interface{}, []interface{}:
					result += "</div>\n"
					result += GenerateReadOnlyJSON(logger, value, fieldPath, indent+1)
				default:
					logger.Debug("Simple value")
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
		logger.Debug("Handle primitive values (should only happen for the root if it's not an object or array)")
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
