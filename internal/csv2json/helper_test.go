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

// TestSetValueInternal tests the setValueInternal function which recursively creates and maps nested dictionaries.
func TestSetValueInternal(t *testing.T) {
	tests := []struct {
		name      string
		hierarchy []string
		value     any
		inside    map[string]any
		want      any
	}{
		{
			name:      "single level",
			hierarchy: []string{"key"},
			value:     "value",
			inside:    map[string]any{},
			want:      "value",
		},
		{
			name:      "two levels",
			hierarchy: []string{"parent", "child"},
			value:     "value",
			inside:    map[string]any{},
			want:      map[string]any{"child": "value"},
		},
		{
			name:      "three levels",
			hierarchy: []string{"level1", "level2", "level3"},
			value:     "value",
			inside:    map[string]any{},
			want:      map[string]any{"level2": map[string]any{"level3": "value"}},
		},
		{
			name:      "existing nested data",
			hierarchy: []string{"parent", "child1", "grandchild"},
			value:     "value",
			inside: map[string]any{
				"parent": map[string]any{
					"child2": "existing",
				},
			},
			want: map[string]any{
				"child1": map[string]any{
					"grandchild": "value",
				},
				"child2": "existing",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setValueInternal(tt.hierarchy, tt.value, tt.inside)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setValueInternal() = %v, want %v", got, tt.want)
			}
		})
	}
}
