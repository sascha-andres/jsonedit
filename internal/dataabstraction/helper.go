package dataabstraction

import "strconv"

// ConvertToType converts the input string `val` to a specified type `t` such as "int", "float", or "bool".
// Returns the converted value as `any` or an error if the conversion fails.
func ConvertToType(t, val string) (any, error) {
	switch t {
	case "int":
		i, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		return i, nil
	case "float":
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, err
		}
		return f, nil
	case "bool":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	return val, nil
}
