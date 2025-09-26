package csv2json

type (

	// Operand represents a structure defining a dynamic or static value for evaluation in conditions.
	Operand struct {

		// Type is whether to use a column or static value (column vs value)
		Type string `json:"type"`

		// Either the column index (or name if working named) or the static value
		Value string `json:"value"`
	}

	// Condition defines a structure for representing a condition with a value and an associated operator.
	Condition struct {

		// Operand1 represents the first operand in a condition, which defines a value and its type for evaluation.
		Operand1 Operand `json:"operand1"`

		// Operand2 represents the second operand in a condition, defining its type and value for evaluation.
		Operand2 Operand `json:"operand2"`

		// Operator specifies the condition operator (e.g., '=', '!=', '>', '<') to be applied for comparison in the mapping configuration.
		Operator string `json:"operator"`

		// Type denotes the data type for the comparison
		Type string `json:"type"`
	}

	// Conditions define a slice of conditions to be evaluated for filtering data
	Conditions []Condition
)
