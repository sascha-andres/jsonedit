package json

import (
	"io"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// InMemoryLoader is a type that loads JSON schema from an in-memory byte slice.
type InMemoryLoader struct {

	// Doc holds the JSON schema data in byte slice format. Used for in-memory loading and manipulation.
	Doc []byte
}

// Load retrieves the schemaFS data and returns it as an io.ReadCloser. An error is returned if the schemaFS cannot be obtained.
func (receiver InMemoryLoader) Load(url string) (any, error) {
	return jsonschema.UnmarshalJSON(io.NopCloser(strings.NewReader(string(receiver.Doc))))
}
