package compare

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

// GetJSONComparison compares two JSON-serializable objects and returns their differences as a string or an error if any occurs.
func GetJSONComparison(a, b interface{}, indent string) (string, error) {
	ref, err := json.MarshalIndent(a, "", indent)
	if err != nil {
		return "", err
	}
	c, err := json.MarshalIndent(b, "", indent)
	if err != nil {
		return "", err
	}
	if diff := cmp.Diff(string(ref), string(c)); diff != "" {
		return fmt.Sprintf("mismatch (-want +got):\n%s", diff), nil
	}
	return "", nil
}
