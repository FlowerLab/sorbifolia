package httprouter

// HandlerFunc defines the handler used by gin middleware as return value.
type HandlerFunc[T any] func(*T)

// HandlersChain defines a HandlerFunc slice.
type HandlersChain[T any] []HandlerFunc[T]

// IRouter defines all router handle interface includes single and group router.
type IRouter[T any] interface {
	IRoutes[T]
	Group(string, ...HandlerFunc[T]) IRouter[T]
}

// IRoutes defines all router handle interface.
type IRoutes[T any] interface {
	Use(...HandlerFunc[T]) IRoutes[T]
	Any(string, ...HandlerFunc[T]) IRoutes[T]

	GET(string, ...HandlerFunc[T]) IRoutes[T]
	POST(string, ...HandlerFunc[T]) IRoutes[T]
	DELETE(string, ...HandlerFunc[T]) IRoutes[T]
	PATCH(string, ...HandlerFunc[T]) IRoutes[T]
	PUT(string, ...HandlerFunc[T]) IRoutes[T]
	OPTIONS(string, ...HandlerFunc[T]) IRoutes[T]
	HEAD(string, ...HandlerFunc[T]) IRoutes[T]
	CONNECT(string, ...HandlerFunc[T]) IRoutes[T]
	TRACE(string, ...HandlerFunc[T]) IRoutes[T]
}
