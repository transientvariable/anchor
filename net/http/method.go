package http

// Enumeration of supported HTTP methods.
const (
	MethodConnect   = "CONNECT"   // RFC 7231, 4.3.6
	MethodDelete    = "DELETE"    // RFC 7231, 4.3.5
	MethodGet       = "GET"       // RFC 7231, 4.3.1
	MethodHead      = "HEAD"      // RFC 7231, 4.3.2
	MethodOptions   = "OPTIONS"   // RFC 7231, 4.3.7
	MethodPatch     = "PATCH"     // RFC 5789
	MethodPost      = "POST"      // RFC 7231, 4.3.3
	MethodPut       = "PUT"       // RFC 7231, 4.3.4
	MethodTrace     = "TRACE"     // RFC 7231, 4.3.8
	MethodCopy      = "COPY"      // RFC 4918 WebDAV, 9.10 - Duplicates a resource.
	MethodLock      = "LOCK"      // RFC 4918 WebDAV, 9.10 - Locks the resource.
	MethodMkcol     = "MKCOL"     // RFC 4918 WebDAV, 9.3  - Creates the collection specified.
	MethodMove      = "MOVE"      // RFC 4918 WebDAV, 9.9  - Moves the resource.
	MethodPropfind  = "PROPFIND"  // RFC 4918 WebDAV, 9.1  - Performs a property find on the server.
	MethodProppatch = "PROPPATCH" // RFC 4918 WebDAV, 9.2  - Sets or removes properties on the server.
	MethodUnlock    = "UNLOCK"    // RFC 4918 WebDAV, 9.11 - Unlocks the resource.

)

// Methods returns a list of RFC 7231/5789 HTTP methods.
func Methods() []string {
	return []string{
		MethodConnect,
		MethodDelete,
		MethodGet,
		MethodHead,
		MethodOptions,
		MethodPatch,
		MethodPost,
		MethodPut,
		MethodTrace,
	}
}

// MethodsWebDAV returns the list of RFC 4918 WebDAV methods.
func MethodsWebDAV() []string {
	return []string{
		MethodDelete,
		MethodGet,
		MethodPost,
		MethodPut,
		MethodCopy,
		MethodLock,
		MethodMkcol,
		MethodMove,
		MethodPropfind,
		MethodProppatch,
		MethodUnlock,
	}
}
