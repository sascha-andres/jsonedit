package csv2json

import (
	"reflect"
	"testing"
)

// TestSetValue tests the setValue function which creates and maps nested dictionaries based on a hierarchy of keys.
func TestSetValue(t *testing.T) {
	tests := []struct {
		name      string
		hierarchy []string
		value     any
		data      map[string]interface{}
		want      map[string]interface{}
	}{
		{
			name:      "simple key",
			hierarchy: []string{"key"},
			value:     "value",
			data:      map[string]interface{}{},
			want:      map[string]interface{}{"key": "value"},
		},
		{
			name:      "nested key",
			hierarchy: []string{"parent", "child"},
			value:     "value",
			data:      map[string]interface{}{},
			want:      map[string]interface{}{"parent": map[string]interface{}{"child": "value"}},
		},
		{
			name:      "deeply nested key",
			hierarchy: []string{"level1", "level2", "level3"},
			value:     "value",
			data:      map[string]interface{}{},
			want:      map[string]interface{}{"level1": map[string]interface{}{"level2": map[string]interface{}{"level3": "value"}}},
		},
		{
			name:      "existing data",
			hierarchy: []string{"parent", "child2"},
			value:     "value2",
			data: map[string]interface{}{
				"parent": map[string]interface{}{
					"child1": "value1",
				},
			},
			want: map[string]interface{}{
				"parent": map[string]interface{}{
					"child1": "value1",
					"child2": "value2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setValue(tt.hierarchy, tt.value, tt.data)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
