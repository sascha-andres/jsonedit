package jsonedit

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

// getJSONComparison compares two JSON-serializable objects and returns their differences as a string or an error if any occurs.
func (app *App) getJSONComparison(a, b interface{}) (string, error) {
	ref, err := json.MarshalIndent(a, "", app.indent)
	if err != nil {
		return "", err
	}
	c, err := json.MarshalIndent(b, "", app.indent)
	if err != nil {
		return "", err
	}
	if diff := cmp.Diff(string(ref), string(c)); diff != "" {
		return fmt.Sprintf("mismatch (-want +got):\n%s", diff), nil
	}
	return "", nil
}
