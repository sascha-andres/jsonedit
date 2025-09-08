package csv2json

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/txtar"
)

var nrData = make(map[string]map[string]string)

type newRecordFuncFactory func(string) NewRecordFunc

// getStoreNewRecordFunc returns a function that processes record and header slices to store data under a specific test key.
// It initializes a data map for the provided test key if it doesn't already exist.
func getStoreNewRecordFunc(test string) NewRecordFunc {
	nrData[test] = make(map[string]string)
	return func(record, header []string) {
		if _, ok := nrData[test]["expect"]; !ok {
			nrData[test]["expect"] = record[0]
		}
	}
}

var nrFunctions = map[string]newRecordFuncFactory{
	"store": getStoreNewRecordFunc,
}

type Expectation struct {
	ConstructError bool `json:"construct_error"`
	Error          bool `json:"error"`
}

type Parameters struct {
	AccessByHeader     bool              `json:"access_by_header"`
	Separator          string            `json:"separator"`
	OutputType         string            `json:"output_type"`
	NestedPropertyName string            `json:"nested_property_name"`
	Logger             bool              `json:"logger"`
	Array              bool              `json:"array"`
	EnvVariables       map[string]string `json:"env_variables"`
	NewRecordFunc      string            `json:"new_record_func"`
	NewRecordExpect    string            `json:"new_record_expect"`
}

// TestMapper tests the Mapper object's functionality using various configurations and expectations from testdata files.
func TestMapper(t *testing.T) {
	files, err := filepath.Glob("testdata/*.txtar")
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		t.Run(filepath.Base(f), func(t *testing.T) {
			content, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}

			ar := txtar.Parse(content)
			var input, config, output, params, expect []byte
			for _, f := range ar.Files {
				switch f.Name {
				case "input.csv":
					input = f.Data
				case "config.json":
					config = f.Data
				case "output":
					output = f.Data
				case "parameters":
					params = f.Data
				case "expectation":
					expect = f.Data
				}
			}

			var parameters Parameters
			if err := json.NewDecoder(strings.NewReader(string(params))).Decode(&parameters); err != nil {
				t.Fatal(err)
			}
			for k, v := range parameters.EnvVariables {
				t.Setenv(k, v)
			}

			var expectation Expectation
			if err := json.NewDecoder(strings.NewReader(string(expect))).Decode(&expectation); err != nil {
				t.Fatal(err)
			}

			options := []OptionFunc{}
			if parameters.Logger {
				options = append(options, WithLogger(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))))
			}
			options = append(options, WithOutputType(parameters.OutputType))
			options = append(options, WithOptions(config))
			options = append(options, WithNamed(parameters.AccessByHeader))
			options = append(options, WithSeparator(parameters.Separator))
			options = append(options, WithNestedPropertyName(parameters.NestedPropertyName))
			options = append(options, WithArray(parameters.Array))

			mapper, err := NewMapper(options...)
			if err != nil && !expectation.ConstructError {
				t.Fatal(err)
			}
			if err != nil && expectation.ConstructError {
				return
			}
			if err == nil && expectation.ConstructError {
				t.Fatal("construction: expected error but got none")
			}

			if parameters.NewRecordFunc != "" {
				f, ok := nrFunctions[parameters.NewRecordFunc]
				if !ok {
					t.Fatal("new record function not found")
				}
				mapper.SetNewRecordFunc(f(t.Name()))
			}

			out, err := mapper.Map(input)
			if err != nil && !expectation.Error {
				t.Fatal(err)
			}
			if err != nil && expectation.Error {
				return
			}
			if err == nil && expectation.Error {
				t.Fatal("run: expected error but got none")
			}

			if parameters.NewRecordFunc != "" {
				if nrData[t.Name()]["expect"] != nrData[t.Name()][parameters.NewRecordExpect] {
					t.Errorf("new record function: expected %s but got %s", nrData[t.Name()][parameters.NewRecordExpect], nrData[t.Name()]["expect"])
				}
			}

			result := string(out)
			expectData := string(output)
			if parameters.Array {
				result, err = prettyPrint(out)
				if err != nil {
					t.Fatal(err)
				}
				expectData, err = prettyPrint(output)
				if err != nil {
					t.Fatal(err)
				}
			}

			if diff := cmp.Diff(strings.TrimSpace(expectData), strings.TrimSpace(result)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
				writeCandidateFile(filepath.Base(f), result)
			}
		})
	}
}

// prettyPrint accepts a JSON byte slice, formats it with indentation, and returns it as a pretty string or an error.
func prettyPrint(in []byte) (string, error) {
	var a any
	err := json.Unmarshal(in, &a)
	if err != nil {
		return "", err
	}
	out, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// writeCandidateFile writes a result string to a new file in the "testdata" directory with a ".candidate" suffix.
func writeCandidateFile(base string, result string) {
	_ = os.WriteFile(filepath.Join("testdata", base+".candidate"), []byte(result), 0644)
}
