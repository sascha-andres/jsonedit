# Mapping Configuration

The mapping configuration is a JSON file that defines how CSV columns are mapped to properties in the output. The configuration has the following structure:

```json
{
  "mapping": {
    "key1": {
      "property": "propertyName",
      "type": "dataType"
    },
    "key2": {
      "property": "nested.property",
      "type": "dataType"
    }
  },
  "calculated": [
    {
      "property": "calculatedProperty",
      "kind": "kindOfCalculation",
      "format": "formatString",
      "type": "dataType",
      "location": "record"
    }
  ],
  "extra_variables": {
    "variable-name": {
      "value": "variable-value"
    }
  }
}
```

Where:
- `key1`, `key2`, etc. are either:
  - Column indices (0, 1, 2, ...) when not using the `-named` flag
  - Column names from the CSV header when using the `-named` flag
- `propertyName` is the name of the property in the output
- `nested.property` demonstrates how to create nested objects using dot notation
- `dataType` is one of:
  - `int` - converts the value to an integer
  - `float` - converts the value to a floating-point number
  - `bool` - converts the value to a boolean
  - `string` (default) - keeps the value as a string

## Multiple Properties from a Single Column

In addition to mapping a column to a single property, you can map a single column to multiple properties using the `properties` array. This is useful when a column contains data that needs to be split into multiple fields or when you want to duplicate a value across multiple properties.

```json
{
  "mapping": {
    "columnKey": {
      "properties": [
        {
          "property": "firstProperty",
          "type": "dataType"
        },
        {
          "property": "nested.secondProperty",
          "type": "dataType"
        }
      ]
    }
  }
}
```

Where:
- `columnKey` is either a column index or name (depending on whether you're using the `-named` flag)
- Each item in the `properties` array defines a separate output property that will receive the value from the same input column
- Each property definition has the same structure as a regular mapping (with `property` and `type` fields)
- You can use dot notation in the `property` field to create nested structures

### Example

Given a CSV with an address column that contains full addresses:

```csv
id,name,address
1,"John Doe","123 Main St, Springfield, IL 62701"
```

You can map the address column to multiple properties:

```json
{
  "mapping": {
    "id": {
      "property": "userId",
      "type": "int"
    },
    "name": {
      "property": "fullName",
      "type": "string"
    },
    "address": {
      "properties": [
        {
          "property": "originalAddress",
          "type": "string"
        },
        {
          "property": "contact.address",
          "type": "string"
        },
        {
          "property": "shipping.address",
          "type": "string"
        }
      ]
    }
  }
}
```

This will produce:

```json
{
  "userId": 1,
  "fullName": "John Doe",
  "originalAddress": "123 Main St, Springfield, IL 62701",
  "contact": {
    "address": "123 Main St, Springfield, IL 62701"
  },
  "shipping": {
    "address": "123 Main St, Springfield, IL 62701"
  }
}
```

## Calculated Fields

Calculated fields allow you to add dynamic values to your output that are not directly derived from the CSV input. These fields are defined in the `calculated` array of the mapping configuration.

Each calculated field has the following properties:
- `property`: The name of the property in the output (supports dot notation for nested objects)
- `kind`: The type of calculation to perform (see below)
- `format`: Additional information for the calculation, varies by kind
- `type`: The data type of the calculated value (`int`, `float`, `bool`, or `string`)
- `location`: Where the calculated field should be applied - either `record` (default) or `document`

### Kinds of Calculated Fields

1. **datetime**: Adds the current date/time formatted according to the format string
   - `format`: A Go time format string (e.g., "2006-01-02" for date, "15:04:05" for time)

2. **application**: Adds application-specific values
   - `format`: Currently only supports "record", which adds the record index (0-based)

3. **environment**: Adds the value of an environment variable
   - `format`: The name of the environment variable to read

4. **extra**: Adds the value of an extra variable defined in the configuration
   - `format`: The name of the extra variable to use
   - Extra variables are defined in the `extra_variables` section of the configuration

5. **mapping**: Maps values from a source field to different output values
   - `format`: Specified as "field:mapping_list" where:
     - `field` is the source field name (when using `-named`) or index
     - `mapping_list` is a comma-separated list of "from=to" pairs
     - A special "default" mapping can be specified for values that don't match any explicit mapping
   - Example: "from-to:a=0,b=1,default=99" maps values from the "from-to" field:
     - "a" becomes 0
     - "b" becomes 1
     - Any other value becomes 99

### Record-Level vs Document-Level Calculated Fields

Calculated fields can be applied at two different levels:

1. **Record-Level Fields** (`location: "record"`): 
   - Applied to each individual record in the output
   - This is the default if no location is specified
   - Always included regardless of output format

2. **Document-Level Fields** (`location: "document"`):
   - Applied to the entire document, not to individual records
   - Only applied when using array output (with `-array` flag) or when using TOML/YAML output formats
   - Typically used for metadata about the entire document
   - Often placed under a top-level property like `_meta`

Document-level calculated fields are useful for adding metadata about the entire dataset, such as:
- Total number of records processed
- Processing timestamp
- Global configuration values

**Note:** Document-level calculated fields are only applied when the output is a single document containing all records (array mode). They are not applied when outputting individual records as separate JSON objects.

### Example

```json
{
  "mapping": {
    "id": {
      "property": "productId",
      "type": "int"
    }
  },
  "calculated": [
    {
      "property": "metadata.recordNumber",
      "kind": "application",
      "format": "record",
      "type": "int",
      "location": "record"
    },
    {
      "property": "metadata.processedDate",
      "kind": "datetime",
      "format": "2006-01-02",
      "type": "string",
      "location": "record"
    },
    {
      "property": "metadata.processedTime",
      "kind": "datetime",
      "format": "15:04:05",
      "type": "string",
      "location": "record"
    },
    {
      "property": "metadata.userHome",
      "kind": "environment",
      "format": "HOME",
      "type": "string",
      "location": "record"
    },
    {
      "property": "metadata.version",
      "kind": "extra",
      "format": "app-version",
      "type": "string",
      "location": "record"
    },
    {
      "property": "_meta.totalRecords",
      "kind": "application",
      "format": "records",
      "type": "int",
      "location": "document"
    },
    {
      "property": "_meta.processedAt",
      "kind": "datetime",
      "format": "2006-01-02 15:04:05",
      "type": "string",
      "location": "document"
    }
  ],
  "extra_variables": {
    "app-version": {
      "value": "1.0.0"
    }
  }
}
```

This configuration would add the following calculated fields:

Record-level fields (added to each record):
- `metadata.recordNumber`: The 0-based index of the record
- `metadata.processedDate`: The current date in YYYY-MM-DD format
- `metadata.processedTime`: The current time in HH:MM:SS format
- `metadata.userHome`: The value of the HOME environment variable
- `metadata.version`: The string "1.0.0" from the extra variable "app-version"

Document-level fields (added to the top-level document when using array output):
- `_meta.totalRecords`: The total number of records processed
- `_meta.processedAt`: The date and time when the document was processed

## Conditional Properties

Conditional properties allow you to include or exclude properties in the output based on specific conditions. This feature is useful when you want to selectively include fields only when certain criteria are met.

### Condition Structure

A condition is defined as an object with the following properties:
- `operator`: The comparison operator to use (e.g., `=`, `!=`, `>`, `<`, `>=`, `<=`)
- `operand1`: The first operand in the comparison
- `operand2`: The second operand in the comparison
- `type`: The data type to use for comparison (`string`, `int`, `float`, or `bool`)

Each operand is an object with:
- `type`: Either `value` for a fixed value or `column` for a value from a CSV column
- `value`: Either a literal value (when `type` is `value`) or a column index/name (when `type` is `column`)

### Example

```json
{
  "mapping": {
    "0": {
      "properties": [
        {
          "property": "id",
          "type": "int"
        },
        {
          "property": "premiumId",
          "type": "int",
          "condition": {
            "operator": "=",
            "operand1": {
              "type": "column",
              "value": "3"
            },
            "operand2": {
              "type": "value",
              "value": "premium"
            },
            "type": "string"
          }
        }
      ]
    }
  }
}
```

In this example:
- The `id` property is always included
- The `premiumId` property is only included when the value in column 3 equals "premium"
- The comparison is done as strings

### Supported Operators

- `=`: Equal to
- `!=`: Not equal to
- `>`: Greater than
- `<`: Less than
- `>=`: Greater than or equal to
- `<=`: Less than or equal to

### Comparison Types

The `type` field in the condition determines how the values are compared:
- `string`: Values are compared as strings (lexicographically)
- `int`: Values are converted to integers before comparison
- `float`: Values are converted to floating-point numbers before comparison
- `bool`: Values are converted to booleans before comparison

### Use Cases

Conditional properties are useful for:
1. Including special fields only for certain types of records
2. Creating different output structures based on data values
3. Implementing business logic in the mapping process
4. Filtering out unwanted or irrelevant data

# Output Behavior

## Without `-array`

When processing multiple CSV rows without the `-array` flag, each row is converted to a separate JSON document and written to the output with newlines between them. This produces a newline-delimited JSON format (NDJSON/JSON Lines), where each line is a valid JSON object, but the file as a whole is not a standard JSON array.

## With `-array`

When using the `-array` flag, all rows are collected into a single array and output as one document.

## Nested Property Output

When using the `-nested-property` flag with the `-array` flag (or when using TOML output which implicitly enables array mode), the output data is nested under the specified property name:

### JSON Output with Nested Property

Will produce:

```json
{
  "items": [
    {
      "property1": 1,
      "property2": {
        "property3": "hello"
      }
    },
    {
      "property1": 2,
      "property2": {
        "property3": "world"
      }
    }
  ]
}
```

### YAML Output with Nested Property

Will produce:

```yaml
items:
  - property1: 1
    property2:
      property3: hello
  - property1: 2
    property2:
      property3: world
```

### TOML Output Format

When using TOML as the output format, the array data is always wrapped in a property. By default, this property is named "data", but you can customize it using the `-nested-property` flag:

Will produce:

```toml
[_meta]
  processedAt = "2023-05-12 15:30:45"
  totalRecords = 2

[[items]]
property1 = 1
property2 = { property3 = "hello" }

[[items]]
property1 = 2
property2 = { property3 = "world" }
```

Note how the document-level calculated fields appear in the `_meta` section at the top of the document, while record-level calculated fields would appear within each record.

# Examples

## Basic Usage

Given the following CSV:

```csv
1,"hello",2.3
```

And this mapping.json:

```json
{
  "mapping": {
    "0": {
      "property": "property1",
      "type": "int"
    },
    "1": {
      "property": "property2.property3",
      "type": "string"
    },
    "2": {
      "property": "property4",
      "type": "float"
    }
  }
}
```

The default output will be:

```json
{
  "property1": 1,
  "property2": {
    "property3": "hello"
  },
  "property4": 2.3
}
```

## Value Mapping Example

Given the following CSV:

```csv
id,status,value
1,"active",10.5
2,"inactive",20.3
3,"pending",15.7
```

And this mapping.json with value mapping:

```json
{
  "mapping": {
    "id": {
      "property": "id",
      "type": "int"
    },
    "status": {
      "property": "originalStatus",
      "type": "string"
    },
    "value": {
      "property": "amount",
      "type": "float"
    }
  },
  "calculated": [
    {
      "property": "statusCode",
      "kind": "mapping",
      "format": "status:active=1,inactive=0,pending=2,default=-1",
      "type": "int",
      "location": "record"
    }
  ]
}
```

Running with the `-named`

Will produce:

```
{"id":1,"originalStatus":"active","amount":10.5,"statusCode":1}
{"id":2,"originalStatus":"inactive","amount":20.3,"statusCode":0}
{"id":3,"originalStatus":"pending","amount":15.7,"statusCode":2}
```

This example demonstrates how to map string status values to numeric codes using the value mapping feature.

## Using Named Columns

Given the following CSV:

```csv
id,name,price
1,"Product A",19.99
2,"Product B",29.99
```

And this mapping.json:

```json
{
  "mapping": {
    "id": {
      "property": "productId",
      "type": "int"
    },
    "name": {
      "property": "productName",
      "type": "string"
    },
    "price": {
      "property": "pricing.retail",
      "type": "float"
    }
  }
}
```

Running with the `-named`

Will produce (in NDJSON format, with each line being a separate JSON document):

```
{"productId":1,"productName":"Product A","pricing":{"retail":19.99}}
{"productId":2,"productName":"Product B","pricing":{"retail":29.99}}
```

Note: The actual output will not be pretty-printed but shown as compact JSON objects, one per line.

## Output as Array in YAML Format

Will produce:

```yaml
- productId: 1
  productName: Product A
  pricing:
    retail: 19.99
- productId: 2
  productName: Product B
  pricing:
    retail: 29.99
```
