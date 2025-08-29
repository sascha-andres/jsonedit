package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/sascha-andres/reuse/flag"

	"github.com/sascha-andres/jsonedit/internal/csv2json"
)

var (
	separator                            = ";"
	outputType                           = "json"
	configurationFile                    string
	inputFilePath                        = "-"
	outputFilePath                       = "-"
	nestedPropertyName                   = ""
	generateArray, accessByHeader, debug bool
)

func init() {
	flag.StringVar(&separator, "separator", separator, "Separator to use for CSV input")
	flag.StringVar(&outputType, "output-type", outputType, "Output type to use for JSON output")
	flag.StringVar(&configurationFile, "configuration-file", configurationFile, "Path to the configuration file")
	flag.StringVar(&nestedPropertyName, "nested-property-name", nestedPropertyName, "Name of the property to use for nested arrays")
	flag.BoolVar(&generateArray, "generate-array", generateArray, "Generate an array for the output")
	flag.BoolVar(&accessByHeader, "access-by-header", accessByHeader, "Access the CSV by header instead of by index")
	flag.StringVar(&inputFilePath, "input-file", inputFilePath, "Path to the input file (- for Stdin)")
	flag.StringVar(&outputFilePath, "output-file", outputFilePath, "Path to the output file (- for Stdout)")
	flag.BoolVar(&debug, "debug", debug, "Enable debug mode")
}

// main serves as the entry point of the application, initializing and executing the primary logic within the program.
func main() {
	flag.Parse()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if configurationFile == "" {
		return errors.New("missing configuration file")
	}

	data, err := os.ReadFile(configurationFile)
	if err != nil {
		return err
	}

	var logger *slog.Logger
	if debug {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}

	mapper, err := csv2json.NewMapper(
		csv2json.WithNamed(accessByHeader),
		csv2json.WithSeparator(separator),
		csv2json.WithOutputType(outputType),
		csv2json.WithNestedPropertyName(nestedPropertyName),
		csv2json.WithArray(generateArray),
		csv2json.WithOptions(data),
		csv2json.WithLogger(logger),
	)
	if err != nil {
		return err
	}

	var reader *os.File
	if inputFilePath == "-" {
		reader = os.Stdin
	} else {
		f, err := os.Open(inputFilePath)
		if err != nil {
			return err
		}
		reader = f
	}

	data, err = io.ReadAll(reader)
	if err != nil {
		return err
	}
	result, err := mapper.Map(data)
	if err != nil {
		return err
	}
	if outputFilePath == "-" {
		fmt.Println(string(result))
	} else {
		err = os.WriteFile(outputFilePath, result, 0640)
	}

	return err
}
