package dijkstra

type edge struct {
	node   string
	weight int
}

type graph struct {
	nodes map[string][]edge
}

func NewGraph() *graph {
	return &graph{nodes: make(map[string][]edge)}
}

func (g *graph) AddEdge(orig, dest string, weight int) {
	g.nodes[orig] = append(g.nodes[orig], edge{node: dest, weight: weight})
	g.nodes[dest] = append(g.nodes[dest], edge{node: orig, weight: weight})
}

func (g *graph) getEdges(node string) []edge {
	return g.nodes[node]
}

func (g *graph) GetPath(orig, dest string) (int, []string) {
	h := newHeap()
	h.push(path{value: 0, nodes: []string{orig}})
	visited := make(map[string]bool)

	for len(*h.values) > 0 {
		// Find the nearest yet to visit node
		p := h.pop()
		node := p.nodes[len(p.nodes)-1]

		if visited[node] {
			continue
		}

		if node == dest {
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