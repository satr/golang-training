# JSON encoding and decoding

## Theory

Go provides built-in support for encoding and decoding JSON data through the `encoding/json` package. This package offers functions `json.Marshal()` to convert Go values into JSON and `json.Unmarshal()` to convert JSON data back into Go values.

The JSON package works with standard Go types such as:
- `string` maps to JSON strings
- `int`, `float64` map to JSON numbers  
- `bool` maps to JSON booleans
- `[]interface{}` maps to JSON arrays
- `map[string]interface{}` maps to JSON objects

To control how struct fields are encoded/decoded, you can use struct tags:

```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
```

This allows customizing field names in the output JSON.

JSON encoding/decoding is particularly useful for HTTP APIs, configuration files, and data interchange between systems.

## Task

Implement the `ParsePerson` function that decodes a JSON string into a Person struct, and `SerializePerson` function that encodes a Person struct into a JSON string. The Person struct should be able to handle both valid and invalid JSON input.

## Run

```bash
make 01_core_12_json_encoding_decoding
```

## Further practice

- Add marshaling and unmarshaling for a more complex nested structure
- Implement custom JSON marshaling with `json.Marshaler` and `json.Unmarshaler` interfaces
- Handle different JSON field name conventions (camelCase, snake_case)