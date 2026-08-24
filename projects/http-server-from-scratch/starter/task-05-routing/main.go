package main

// Response is the HTTP response type used by handlers.
type Response struct {
	Status int
	Body   string
}

// HandlerFunc is a function that handles an HTTP request.
type HandlerFunc func(method, path string, headers map[string]string) Response

// Mux routes incoming requests to registered handlers.
type Mux struct {
	routes map[string]HandlerFunc // key: "METHOD /path"
}

// NewMux creates a new Mux with an empty route table.
func NewMux() *Mux {
	return &Mux{routes: make(map[string]HandlerFunc)}
}

// Handle registers fn as the handler for "METHOD /path".
func (m *Mux) Handle(method, path string, fn HandlerFunc) {
	// TODO: store fn in m.routes with key method + " " + path
}

// Dispatch finds and calls the handler for method+path, or returns a 404 Response.
func (m *Mux) Dispatch(method, path string, headers map[string]string) Response {
	// TODO: look up method + " " + path in m.routes.
	// If found, call and return the handler result.
	// Otherwise return Response{Status: 404}.
	return Response{Status: 404}
}

func main() {}
