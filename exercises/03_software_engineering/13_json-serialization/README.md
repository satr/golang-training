# JSON Serialization

## Theory

Go's `encoding/json` package uses struct tags to control field naming and behavior. Custom `MarshalJSON` and `UnmarshalJSON` methods let you control serialization format, validate on input, and hide internal fields.

## Task

Implement a JSON-serializable struct with custom marshaling, struct tags for field naming, and a custom UnmarshalJSON that validates input.

## Run

```bash
make 03_software_engineering_13_json_serialization
```

## Further practice

Add omitempty behavior, custom type serializers for UUIDs and durations.
