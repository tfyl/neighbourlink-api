package dijkstra

// structure for each edge
type edge struct {
	node   string
	weight int
}

// structure for graph
type graph struct {
	nodes map[string][]edge
}

// creates new graph by instantiating the struct/object
func NewGraph() *graph {
	return &graph{nodes: make(map[string][]edge)}
}

// public method for the graph struct/object ( called with object.AddEdge() )
// adds edge
func (g *graph) AddEdge(orig, dest string, weight int) {
	// adds the connections for origin
	g.nodes[orig] = append(g.nodes[orig], edge{node: dest, weight: weight})
	// adds connection for the destination node
	g.nodes[dest] = append(g.nodes[dest], edge{node: orig, weight: weight})
}

// private method for the graph struct/object ( called with object.getEdges() )
// returns te edges for a given node
func (g *graph) getEdges(node string) []edge {
	return g.nodes[node]
}

// public method for the graph struct/object ( called with object.GetPath() )
// gets shortest/safest path between two nodes in the graph
func (g *graph) GetPath(orig, dest string) (int, []string) {
	// instantiates the heap
	h := newHeap()
	// pushes the origin node to the heap
	h.push(path{value: 0, nodes: []string{orig}})
	// map/dictionary for all visited nodes
	visited := make(map[string]bool)

	for len(*h.values) > 0 {
		// Find the nearest yet to visit node
		p := h.pop()
		node := p.nodes[len(p.nodes)-1]

		if visited[node] {
			continue
		}

		if node == dest {
			// return value as the destination has been reached
			return p.value, p.nodes
		}

		for _, e := range g.getEdges(node) {
			if !visited[e.node] {
				// We calculate the total spent so far plus the cost and the path of getting here
				h.push(path{value: p.value + e.weight, nodes: append([]string{}, append(p.nodes, e.node)...)})
			}
		}

		visited[node] = true
	}

	return 0, nil
}