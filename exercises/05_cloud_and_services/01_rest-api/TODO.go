//go:build todo

package restapi

// Item represents a domain entity stored and returned via JSON.
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Store abstracts persistence for items.
type Store interface {
	Get(ctx context.Context, id string) (Item, error)
	List(ctx context.Context) ([]Item, error)
	Create(ctx context.Context, item Item) (Item, error)
}

// NotFoundError indicates that a requested item could not be found.
type NotFoundError struct {
	ID string
}

func (e *NotFoundError) Error() string {
	// TODO: Return a message such as "not found: <id>".
	return ""
}

// Handler serves HTTP requests against a Store.
type Handler struct {
	store  Store
	nextID uint64
	mu     sync.Mutex
}

// NewHandler returns a Handler that uses the given store.
func NewHandler(store Store) *Handler {
	// TODO: Return a non-nil *Handler.
	return nil
}

// ServeHTTP dispatches GET and POST requests to /items.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Route by method+path:
	//   GET  /items           -> 200, JSON array of items
	//   GET  /items/{id}      -> 200, single item JSON; 404 on NotFoundError
	//   POST /items           -> 201, created item JSON (auto ID)
	//   everything else       -> 405
}
