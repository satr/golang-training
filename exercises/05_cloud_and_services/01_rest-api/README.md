# REST API design

## Theory

A well-structured HTTP handler wraps business logic behind an interface.
Calling the store through an interface decouples the handler from any specific
database, cache, or remote API. Tests can inject an in-memory store, and
production code can wire up PostgreSQL, S3, or a gRPC client.

The `ServeHTTP` method inspects `Method` and `URL.Path` to dispatch work:

- `GET /items` → list items
- `GET /items/{id}` → get a single item by ID
- `POST /items` → create a new item
- anything else → `405 Method Not Allowed`

JSON responses use the standard `encoding/json` package. Error responses use
`http.Error`.

## Task

Copy `solution.go` from `TODO.go`, then implement:

1. `NotFoundError.Error` — return a message that includes the item ID.
2. `NewHandler` — return a non-nil `*Handler`.
3. `ServeHTTP` — dispatch by method and path. Use the store for list/get/create.
4. Auto-generate an ID for every `POST /items` (use the `nextID` field).

## Run

```bash
make 05_cloud_and_services_01_rest-api
```

## Further practice

Add pagination to `GET /items` (e.g. `?limit=10&offset=0`). Add a `PUT /items/{id}` endpoint.
