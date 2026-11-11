# HTTP Handlers

## Theory

HTTP handlers in Go are functions that read `*http.Request` and write to
`http.ResponseWriter`. Use proper status codes (200 OK, 201 Created, 400 Bad
Request), "application/json" content-type headers, and json encoding for APIs.

## Task

Implement `CreateItem` and `GetItems` HTTP handlers that work with an
in-memory store.

- `CreateItem` accepts POST requests with a JSON body `{name, quantity}`.
  Name is required; return 400 if missing. Return 400 on invalid JSON.
- `GetItems` returns all stored items as a JSON array.

## Run

```
make 03_software_engineering_04_http_handlers
```

## Further practice

- Add a DELETE /item handler.
- Add item validation beyond name check.
