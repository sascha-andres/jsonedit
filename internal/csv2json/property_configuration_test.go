package csv2json

import (
	"io"
	"log/slog"
	"testing"
)

func TestPropertyConfiguration_Applies(t *testing.T) {
	// Create a logger that discards all output
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name   string
		pc     PropertyConfiguration
		named  bool
		record []string
		header []string
		want   bool
	}{
		{
			name: "nil condition",
			pc: PropertyConfiguration{
				Property:  "test",
				Type:      "string",
				Condition: nil,
			},
			named:  false,
			record: []string{"value1", "value2"},
			header: []string{"header1", "header2"},
			want:   true,
		},
		{
			name: "int condition equals true",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "int",
				Condition: &Condition{
					Type:     "int",
					Operator: "=",
					Operand1: Operand{Type: "value", Value: "10"},
					Operand2: Operand{Type: "value", Value: "10"},
				},
			},
			named:  false,
			record: []string{},
			header: []string{},
			want:   true,
		},
		{
			name: "int condition equals false",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "int",
				Condition: &Condition{
					Type:     "int",
					Operator: "=",
					Operand1: Operand{Type: "value", Value: "10"},
					Operand2: Operand{Type: "value", Value: "20"},
				},
			},
			named:  false,
			record: []string{},
			header: []string{},
			want:   false,
		},
		{
			name: "int condition not equals true",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "int",
				Condition: &Condition{
					Type:     "int",
					Operator: "!=",
					Operand1: Operand{Type: "value", Value: "10"},
					Operand2: Operand{Type: "value", Value: "20"},
				},
			},
			named:  false,
			record: []string{},
			header: []string{},
			want:   true,
		},
		{
			name: "int condition greater than true",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "int",
				Condition: &Condition{
					Type:     "int",
					Operator: ">",
					Operand1: Operand{Type: "value", Value: "20"},
					Operand2: Operand{Type: "value", Value: "10"},
				},
			},
			named:  false,
			record: []string{},
			header: []string{},
			want:   true,
		},
		{
			name: "int condition less than true",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "int",
				Condition: &Condition{
					Type:     "int",
					Operator: "<",
					Operand1: Operand{Type: "value", Value: "10"},
					Operand2: Operand{Type: "value", Value: "20"},
				},
			},
			named:  false,
			record: []string{},
			header: []string{},
			want:   true,
		},
		{
			name: "int condition with column reference",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "int",
				Condition: &Condition{
					Type:     "int",
					Operator: "=",
					Operand1: Operand{Type: "column", Value: "0"},
					Operand2: Operand{Type: "value", Value: "10"},
				},
			},
			named:  false,
			record: []string{"10", "20"},
			header: []string{"id", "value"},
			want:   true,
		},
		{
			name: "float condition equals true",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "float",
				Condition: &Condition{
					Type:     "float",
					Operator: "=",
					Operand1: Operand{Type: "value", Value: "10.5"},
					Operand2: Operand{Type: "value", Value: "10.5"},
				},
			},
			named:  false,
			record: []string{},
			header: []string{},
			want:   true,
		},
		{
			name: "float condition not equals true",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "float",
				Condition: &Condition{
					Type:     "float",
					Operator: "!=",
					Operand1: Operand{Type: "value", Value: "10.5"},
					Operand2: Operand{Type: "value", Value: "20.5"},
				},
			},
			named:  false,
			record: []string{},
			header: []string{},
			want:   true,
		},
		{
			name: "float condition greater than true",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "float",
				Condition: &Condition{
					Type:     "float",
					Operator: ">",
					Operand1: Operand{Type: "value", Value: "20.5"},
					Operand2: Operand{Type: "value", Value: "10.5"},
				},
			},
			named:  false,
			record: []string{},
			header: []string{},
			want:   true,
		},
		{
			name: "float condition with column reference",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "float",
				Condition: &Condition{
					Type:     "float",
					Operator: "=",
					Operand1: Operand{Type: "column", Value: "0"},
					Operand2: Operand{Type: "value", Value: "10.5"},
				},
			},
			named:  false,
			record: []string{"10.5", "20.5"},
			header: []string{"price", "quantity"},
			want:   true,
		},
		{
			name: "string condition equals true",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "string",
				Condition: &Condition{
					Type:     "string",
					Operator: "=",
					Operand1: Operand{Type: "value", Value: "hello"},
					Operand2: Operand{Type: "value", Value: "hello"},
				},
			},
			named:  false,
			record: []string{},
			header: []string{},
			want:   true,
		},
		{
			name: "string condition not equals true",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "string",
				Condition: &Condition{
					Type:     "string",
					Operator: "!=",
					Operand1: Operand{Type: "value", Value: "hello"},
					Operand2: Operand{Type: "value", Value: "world"},
				},
			},
			named:  false,
			record: []string{},
			header: []string{},
			want:   true,
		},
		{
			name: "string condition with column reference",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "string",
				Condition: &Condition{
					Type:     "string",
					Operator: "=",
					Operand1: Operand{Type: "column", Value: "0"},
					Operand2: Operand{Type: "value", Value: "hello"},
				},
			},
			named:  false,
			record: []string{"hello", "world"},
			header: []string{"greeting", "target"},
			want:   true,
		},
		{
			name: "bool condition equals true",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "bool",
				Condition: &Condition{
					Type:     "bool",
					Operator: "=",
					Operand1: Operand{Type: "value", Value: "true"},
					Operand2: Operand{Type: "value", Value: "ignored"}, // For bool, only Operand1 is used
				},
			},
			named:  false,
			record: []string{},
			header: []string{},
			want:   true,
		},
		{
			name: "bool condition not equals true",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "bool",
				Condition: &Condition{
					Type:     "bool",
					Operator: "!=",
					Operand1: Operand{Type: "value", Value: "false"},
					Operand2: Operand{Type: "value", Value: "ignored"}, // For bool, only Operand1 is used
				},
			},
			named:  false,
			record: []string{},
			header: []string{},
			want:   true,
		},
		{
			name: "bool condition with column reference",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "bool",
				Condition: &Condition{
					Type:     "bool",
					Operator: "=",
					Operand1: Operand{Type: "column", Value: "0"},
					Operand2: Operand{Type: "value", Value: "ignored"}, // For bool, only Operand1 is used
				},
			},
			named:  false,
			record: []string{"true", "other"},
			header: []string{"active", "name"},
			want:   true,
		},
		{
			name: "bool condition with unsupported operator",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "bool",
				Condition: &Condition{
					Type:     "bool",
					Operator: ">", // Not supported for bool
					Operand1: Operand{Type: "value", Value: "true"},
					Operand2: Operand{Type: "value", Value: "ignored"},
				},
			},
			named:  false,
			record: []string{},
			header: []string{},
			want:   false,
		},
		{
			name: "unsupported condition type",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "unknown",
				Condition: &Condition{
					Type:     "unknown",
					Operator: "=",
					Operand1: Operand{Type: "value", Value: "value"},
					Operand2: Operand{Type: "value", Value: "value"},
				},
			},
			named:  false,
			record: []string{},
			header: []string{},
			want:   false,
		},
		{
			name: "named column reference",
			pc: PropertyConfiguration{
				Property: "test",
				Type:     "string",
				Condition: &Condition{
					Type:     "string",
					Operator: "=",
					Operand1: Operand{Type: "column", Value: "name"},
					Operand2: Operand{Type: "value", Value: "John"},
				},
			},
			named:  true,
			record: []string{"John", "30"},
			header: []string{"name", "age"},
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recordInfo := &RecordWithInformation{
				Record: tt.record,
				Header: tt.header,
			}
			if tt.named {
				// build header index cache
				headerIndex := make(map[string]int, len(tt.header))
				for i, h := range tt.header {
					headerIndex[h] = i
				}
				recordInfo.HeaderIndex = headerIndex
			}
			got := tt.pc.Applies(logger, recordInfo)
			if got != tt.want {
				t.Errorf("PropertyConfiguration.Applies() = %v, want %v", got, tt.want)
			}
		})
	}
}
