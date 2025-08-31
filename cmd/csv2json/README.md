# csv2json

A command-line utility for converting CSV data to JSON format with advanced configuration options.

## Overview

The csv2json utility transforms CSV data into JSON, YAML, or TOML formats with flexible mapping options. It allows you to:

- Convert CSV files to structured data formats
- Configure type conversions for each column
- Create nested JSON structures
- Add calculated fields
- Generate arrays or objects
- Read from files or standard input
- Write to files or standard output

## Installation

### From Source

```bash
go install github.com/sascha-andres/jsonedit/cmd/csv2json@latest
```

### From Binary Releases

Download the appropriate binary for your platform from the [releases page](https://github.com/sascha-andres/jsonedit/releases).

## Usage

```bash
csv2json --configuration-file=config.json --input-file=data.csv --output-file=result.json
```

### Basic Example

Convert a CSV file to JSON:

```bash
csv2json --configuration-file=config.json --input-file=input.csv --output-file=output.json
```

Read from stdin and write to stdout:

```bash
cat input.csv | csv2json --configuration-file=config.json > output.json
```

## Command-Line Options

| Option | Default | Description |
|--------|---------|-------------|
| `--separator` | `;` | Separator character used in the CSV input |
| `--output-type` | `json` | Output format type (`json`, `yaml`, or `toml`) |
| `--configuration-file` | (required) | Path to the configuration file |
| `--nested-property-name` | (empty) | Name of the property to use for nested arrays (for TOML output) |
| `--generate-array` | `false` | Generate an array for the output instead of individual objects |
| `--access-by-header` | `false` | Access CSV columns by header name instead of by index |
| `--input-file` | `-` (stdin) | Path to the input file (`-` for stdin) |
| `--output-file` | `-` (stdout) | Path to the output file (`-` for stdout) |
| `--debug` | `false` | Enable debug mode with detailed logging to stdout |

## Configuration File Format

The configuration file defines how CSV columns are mapped to the output format. It supports JSON, YAML, or TOML formats.

### Example Configuration

```json
{
  "mapping": {
    "0": {
      "property": "id",
      "type": "int"
    },
    "1": {
      "property": "name",
      "type": "string"
    },
    "2": {
      "properties": [
        {
          "property": "address.street",
          "type": "string"
        },
        {
          "property": "address.city",
          "type": "string"
        }
      ]
    }
  },
  "calculated": [
    {
      "property": "timestamp",
      "kind": "datetime",
      "format": "2006-01-02T15:04:05Z07:00",
      "type": "string",
      "location": "document"
    }
  ],
  "extra_variables": {
    "version": {
      "value": "1.0"
    }
  }
}
```

### Configuration Structure

#### Mapping

The `mapping` section maps CSV columns to output properties:

- Keys can be column indices (when `--access-by-header=false`) or header names (when `--access-by-header=true`)
- Each mapping can have either:
  - A single `property` and `type` for direct mapping
  - Multiple `properties` for mapping one column to multiple output fields

#### Conditional Mapping

Properties can be conditionally included in the output based on specified criteria:

- Add a `condition` object to a property mapping to make it conditional
- The condition compares two operands using an operator
- The property will only be included in the output when the condition is true

#### Supported Types

- `string` - Text values
- `int` - Integer values
- `float` - Floating-point values
- `bool` - Boolean values (true/false)

#### Calculated Fields

The `calculated` section defines dynamic fields:

- `property`: The name of the property in the output
- `kind`: The type of calculation (`datetime` or `application`)
- `format`: Format string (for datetime) or value to acquire
- `type`: Output data type
- `location`: Where to apply the field (`document` or `record`)

#### Extra Variables

The `extra_variables` section defines static variables to include in the output.

## Examples

### Simple CSV to JSON Conversion

For a CSV file `users.csv`:
```
1;John Doe;john@example.com
2;Jane Smith;jane@example.com
```

With configuration `config.json`:
```json
{
  "mapping": {
    "0": {
      "property": "id",
      "type": "int"
    },
    "1": {
      "property": "name",
      "type": "string"
    },
    "2": {
      "property": "email",
      "type": "string"
    }
  }
}
```

Command:
```bash
csv2json --configuration-file=config.json --input-file=users.csv --output-file=users.json
```

Output (in JSON Lines format, with each line being a valid JSON object):
```jsonl
{"id":1,"name":"John Doe","email":"john@example.com"}
{"id":2,"name":"Jane Smith","email":"jane@example.com"}
```

### Generating an Array

Using the same example but with the `--generate-array` flag:

```bash
csv2json --configuration-file=config.json --input-file=users.csv --output-file=users.json --generate-array
```

Output:
```json
[
  {"id":1,"name":"John Doe","email":"john@example.com"},
  {"id":2,"name":"Jane Smith","email":"jane@example.com"}
]
```

## Error Handling

The utility will exit with an error message if:
- The configuration file is missing or invalid
- The input file cannot be read
- The output file cannot be written
- The CSV data is malformed
- Type conversions fail
