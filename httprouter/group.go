package httprouter

type Group[T any] struct {
	Handler  Handlers[T]
	BasePath string
	route    *Router[T]
}

func (g *Group[T]) GET(s string, h ...Handler[T]) IRoutes[T]     { return g.addRoute(GET, s, h) }
func (g *Group[T]) POST(s string, h ...Handler[T]) IRoutes[T]    { return g.addRoute(POST, s, h) }
func (g *Group[T]) DELETE(s string, h ...Handler[T]) IRoutes[T]  { return g.addRoute(DELETE, s, h) }
func (g *Group[T]) PATCH(s string, h ...Handler[T]) IRoutes[T]   { return g.addRoute(PATCH, s, h) }
func (g *Group[T]) PUT(s string, h ...Handler[T]) IRoutes[T]     { return g.addRoute(PUT, s, h) }
func (g *Group[T]) OPTIONS(s string, h ...Handler[T]) IRoutes[T] { return g.addRoute(OPTIONS, s, h) }
func (g *Group[T]) HEAD(s string, h ...Handler[T]) IRoutes[T]    { return g.addRoute(HEAD, s, h) }
func (g *Group[T]) CONNECT(s string, h ...Handler[T]) IRoutes[T] { return g.addRoute(CONNECT, s, h) }
func (g *Group[T]) TRACE(s string, h ...Handler[T]) IRoutes[T]   { return g.addRoute(TRACE, s, h) }
func (g *Group[T]) Use(h ...Handler[T]) IRoutes[T]               { g.Handler = append(g.Handler, h...); return g }
func (g *Group[T]) Any(s string, h ...Handler[T]) IRoutes[T] {
	for _, v := range methods {
		g.addRoute(v, s, h)
	}
	return g
}

func (g *Group[T]) Group(path string, handler ...Handler[T]) IRouter[T] {
	return &Group[T]{
		Handler:  g.combineHandlers(handler),
		BasePath: g.BasePath + path,
		route:    g.route,
	}
}

func (g *Group[T]) addRoute(method Method, path string, handler Handlers[T]) *Group[T] {
	g.route.AddRoute(method, g.BasePath+path, g.combineHandlers(handler))
	return g
}

func (g *Group[T]) combineHandlers(handler Handlers[T]) Handlers[T] {
	return append(g.Handler, handler...)
}
