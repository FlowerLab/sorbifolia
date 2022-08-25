package httprouter

// Handler defines the handler used by gin middleware as return value.
type Handler[T any] func(*T)

// Handlers defines a Handler slice.
type Handlers[T any] []Handler[T]

// IRouter defines all router handle interface includes single and group router.
type IRouter[T any] interface {
	IRoutes[T]
	Group(string, ...Handler[T]) IRouter[T]
}

// IRoutes defines all router handle interface.
type IRoutes[T any] interface {
	Use(...Handler[T]) IRoutes[T]
	Any(string, ...Handler[T]) IRoutes[T]

	GET(string, ...Handler[T]) IRoutes[T]
	POST(string, ...Handler[T]) IRoutes[T]
	DELETE(string, ...Handler[T]) IRoutes[T]
	PATCH(string, ...Handler[T]) IRoutes[T]
	PUT(string, ...Handler[T]) IRoutes[T]
	OPTIONS(string, ...Handler[T]) IRoutes[T]
	HEAD(string, ...Handler[T]) IRoutes[T]
	CONNECT(string, ...Handler[T]) IRoutes[T]
	TRACE(string, ...Handler[T]) IRoutes[T]
}
