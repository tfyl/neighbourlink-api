package types

// class / struct that defines how the shortest path
// it is used when the shortest path from a graph needs to returned
// or received

type ShortestPath struct {
	Cost int
	Path []string
}