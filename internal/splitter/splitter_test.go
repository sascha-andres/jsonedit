package splitter

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

type Expectations struct {
	IsArray                      bool   `json:"is_array"`
	IsArrayErr                   bool   `json:"is_array_err"`
	ArrayPath                    string `json:"array_path"`
	AllArrayElementsAreValues    bool   `json:"all_array_elements_are_values"`
	AllArrayElementsAreValuesErr bool   `json:"all_array_elements_are_values_err"`
	ExtractError                 bool   `json:"extract_error"`
	WithEmptyGroups              bool   `json:"with_empty_groups"`
	SplitError                   bool   `json:"split_error"`
}

// TestAllArrayElementsAreValues verifies that all elements in arrays across multiple test cases are primitive value types.
func TestAllArrayElementsAreValues(t *testing.T) {
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
			var input, expect []byte
			for _, f := range ar.Files {
				switch f.Name {
				case "input.json":
					input = f.Data
				case "expectations.json":
					expect = f.Data
				}
			}

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			var e Expectations
			if err := json.Unmarshal(expect, &e); err != nil {
				t.Fatalf("error getting expectations: %s", err)
			}

			s, err := NewSplitter(WithLogger(logger), WithArrayPath(e.ArrayPath))
			if err != nil {
				t.Fatalf("error creating splitter: %s", err)
			}
			var toSplit any
			err = json.Unmarshal(input, &toSplit)
			if err != nil && !e.IsArrayErr {
				t.Fatalf("did not expect error while unmarshalling: %s", err)
			}
			if err == nil && e.IsArrayErr {
				t.Fatalf("expected error while unmarshalling")
			}

			data, err := s.extractDataToSplit(toSplit)
			if err != nil {
				if e.ExtractError {
					return
				}
				t.Fatalf("error extracting data to split: %s", err)
			}
			if e.ExtractError {
				t.Fatalf("expected error extracting data to split")
			}
			allAreValues := s.allArrayElementsAreValues(data)
			if allAreValues != e.AllArrayElementsAreValues {
				t.Errorf("expected all array elements to be values to be %t but got %t", e.AllArrayElementsAreValues, allAreValues)
			}
		})
	}
}

// TestIsArray verifies the functionality of checking whether input data is a JSON array or not, based on expectations.
func TestIsArray(t *testing.T) {
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
			var input, expect []byte
			for _, f := range ar.Files {
				switch f.Name {
				case "input.json":
					input = f.Data
				case "expectations.json":
					expect = f.Data
				}
			}

			var e Expectations
			if err := json.Unmarshal(expect, &e); err != nil {
				t.Fatalf("error getting expectations: %s", err)
			}

			var toSplit any
			err = json.Unmarshal(input, &toSplit)
			if err != nil && !e.IsArrayErr {
				t.Fatalf("did not expect error while unmarshalling: %s", err)
			}
			if err == nil && e.IsArrayErr {
				t.Fatalf("expected error while unmarshalling")
			}

			s := Splitter{}
			isJsonArray := s.isJSONArray(toSplit)
			if isJsonArray != e.IsArray {
				t.Errorf("expected is array to be %t but got %t", e.IsArray, isJsonArray)
			}
		})
	}
}

func TestSplit(t *testing.T) {
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
			var input, expect, config, outcome []byte
			for _, f := range ar.Files {
				switch f.Name {
				case "input.json":
					input = f.Data
				case "expectations.json":
					expect = f.Data
				case "config.json":
					config = f.Data
				case "outcome.json":
					outcome = f.Data
				}
			}

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

			var e Expectations
			if err := json.Unmarshal(expect, &e); err != nil {
				t.Fatalf("error getting expectations: %s", err)
			}

			var cfg Configuration
			err = json.Unmarshal(config, &cfg)
			if err != nil {
				t.Fatalf("error getting config: %s", err)
			}

			s, err := NewSplitter(
				WithLogger(logger),
				WithArrayPath(e.ArrayPath),
				WithConfiguration(cfg),
				WithOutputEmptyGroups(e.WithEmptyGroups),
			)
			if err != nil {
				t.Fatalf("error creating splitter: %s", err)
			}
			var toSplit any
			err = json.Unmarshal(input, &toSplit)
			if err != nil && !e.IsArrayErr {
				t.Fatalf("did not expect error while unmarshalling: %s", err)
			}
			if err == nil && e.IsArrayErr {
				t.Fatalf("expected error while unmarshalling")
			}

			data, err := s.extractDataToSplit(toSplit)
			if err != nil {
				if e.ExtractError {
					return
				}
				t.Fatalf("error extracting data to split: %s", err)
			}
			if e.ExtractError {
				t.Fatalf("expected error extracting data to split")
			}
			allAreValues := s.allArrayElementsAreValues(data)
			if allAreValues != e.AllArrayElementsAreValues {
				t.Errorf("expected all array elements to be values to be %t but got %t", e.AllArrayElementsAreValues, allAreValues)
			}

			result, err := s.Split(input)
			if e.SplitError && err == nil {
				t.Fatalf("expected error splitting")
			}
			if !e.SplitError && err != nil {
				t.Fatalf("unexpected error splitting: %s", err)
			}
			if e.SplitError && err != nil {
				t.Logf("expected error splitting: %s", err)
				return
			}
			res, err := json.MarshalIndent(result, "", "  ")
			if diff := cmp.Diff(strings.TrimSpace(string(outcome)), strings.TrimSpace(string(res))); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
				writeCandidateFile(filepath.Base(f), string(res))
			}
		})
	}
}

// writeCandidateFile writes a result string to a new file in the "testdata" directory with a ".candidate" suffix.
func writeCandidateFile(base string, result string) {
	_ = os.WriteFile(filepath.Join("testdata", base+".candidate"), []byte(result), 0644)
}
