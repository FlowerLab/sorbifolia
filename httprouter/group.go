package httprouter

type Group[T any] struct {
	Handler  HandlersChain[T]
	BasePath string
	route    *Router[T]
}

func (g *Group[T]) GET(s string, h ...HandlerFunc[T]) IRoutes[T]     { return g.addRoute(GET, s, h) }
func (g *Group[T]) POST(s string, h ...HandlerFunc[T]) IRoutes[T]    { return g.addRoute(GET, s, h) }
func (g *Group[T]) DELETE(s string, h ...HandlerFunc[T]) IRoutes[T]  { return g.addRoute(GET, s, h) }
func (g *Group[T]) PATCH(s string, h ...HandlerFunc[T]) IRoutes[T]   { return g.addRoute(GET, s, h) }
func (g *Group[T]) PUT(s string, h ...HandlerFunc[T]) IRoutes[T]     { return g.addRoute(GET, s, h) }
func (g *Group[T]) OPTIONS(s string, h ...HandlerFunc[T]) IRoutes[T] { return g.addRoute(GET, s, h) }
func (g *Group[T]) HEAD(s string, h ...HandlerFunc[T]) IRoutes[T]    { return g.addRoute(GET, s, h) }
func (g *Group[T]) CONNECT(s string, h ...HandlerFunc[T]) IRoutes[T] { return g.addRoute(GET, s, h) }
func (g *Group[T]) TRACE(s string, h ...HandlerFunc[T]) IRoutes[T]   { return g.addRoute(GET, s, h) }
func (g *Group[T]) Use(h ...HandlerFunc[T]) IRoutes[T]               { g.Handler = append(g.Handler, h...); return g }
func (g *Group[T]) Any(s string, h ...HandlerFunc[T]) IRoutes[T] {
	for _, v := range []Method{
		GET, HEAD, POST, PUT, PATCH, DELETE, CONNECT, OPTIONS, TRACE,
	} {
		g.addRoute(v, s, h)
	}
	return g
}

func (g *Group[T]) Group(path string, handler ...HandlerFunc[T]) IRouter[T] {
	return &Group[T]{
		Handler:  g.combineHandlers(handler),
		BasePath: g.BasePath + path,
		route:    g.route,
	}
}

func (g *Group[T]) addRoute(method Method, path string, handler HandlersChain[T]) *Group[T] {
	g.route.AddRoute(method, g.BasePath+path, g.combineHandlers(handler))
	return g
}

func (g *Group[T]) combineHandlers(handler HandlersChain[T]) HandlersChain[T] {
	return append(g.Handler, handler...)
}
