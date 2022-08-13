package datastructures

import (
	"errors"
)

var (
	ErrVertexNotFound = errors.New("vertex not found")
	ErrSelfLoop       = errors.New("self loops not permitted")
	ErrParallelEdge   = errors.New("parallel edges are not permitted")
)

// Graph is a mutable, non-persistent undirected graph.
// Parallel edges and self-loops are not permitted.
// Additional description: https://en.wikipedia.org/wiki/Graph_(discrete_mathematics)#Simple_graph
type Graph[T comparable] struct {
	adjacencyList map[T]map[T]struct{}
	v, e          int
}

func (g *Graph[T]) V() int { return g.v }
func (g *Graph[T]) E() int { return g.e }

// AddEdge will create an edge between vertices v and w
func (g *Graph[T]) AddEdge(v, w T) error {
	if v == w {
		return ErrSelfLoop
	}

	g.addVertex(v)
	g.addVertex(w)

	if _, ok := g.adjacencyList[v][w]; ok {
		return ErrParallelEdge
	}

	g.adjacencyList[v][w] = struct{}{}
	g.adjacencyList[w][v] = struct{}{}
	g.e++
	return nil
}

// Adj returns the list of all vertices connected to v
func (g *Graph[T]) Adj(v T) ([]T, error) {
	deg, err := g.Degree(v)
	if err != nil {
		return nil, ErrVertexNotFound
	}

	adj := make([]T, deg)
	i := 0
	for key := range g.adjacencyList[v] {
		adj[i] = key
		i++
	}
	return adj, nil
}

// Degree returns the number of vertices connected to v
func (g *Graph[T]) Degree(v T) (int, error) {
	val, ok := g.adjacencyList[v]
	if !ok {
		return 0, ErrVertexNotFound
	}
	return len(val), nil
}

func (g *Graph[T]) addVertex(v T) {
	if _, ok := g.adjacencyList[v]; !ok {
		g.adjacencyList[v] = make(map[T]struct{})
		g.v++
	}
}

// NewGraph creates and returns a Graph
func NewGraph[T comparable]() *Graph[T] {
	return &Graph[T]{
		adjacencyList: make(map[T]map[T]struct{}),
		v:             0,
		e:             0,
	}
}
