package http

import (
	"arena"
)

type ClientArena struct {
	A *arena.Arena
}

func (c *ClientArena) Do(req *Request, resp *Response) {

}

func (c *ClientArena) NewRequest() *Request {
	return arena.New[Request](c.A)
}

func (c *ClientArena) NewResponse() *Response {
	return arena.New[Response](c.A)
}
