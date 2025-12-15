package main

import (
	"fmt"
	"strings"

	"github.com/mbordner/aoc2025/common/files"
)

func main() {
	nodes := getNodes("../data.txt")
	/*
		paths := getPaths(nodes.GetOrCreate("svr"), nodes.GetOrCreate("out"))
		fmt.Println(len(paths))
		fmt.Println(paths)

		fmt.Println(countPaths(nodes.GetOrCreate("svr"), nodes.GetOrCreate("out")))
	*/
	//wg := &sync.WaitGroup{}
	//
	for _, path := range [][]string{{"svr", "dac"}, {"fft", "out"}, {"dac", "fft"}, {"svr", "fft"}, {"fft", "dac"}, {"dac", "out"}} {
		s := nodes.GetOrCreate(path[0])
		d := nodes.GetOrCreate(path[1])
		//wg.Add(1)
		//go func() {
		fmt.Println(fmt.Sprintf("%s->%s", s.ID(), d.ID()), countPaths(s, d))
		//}()
	}
	//wg.Wait()
}

type NodePaths []NodePath

func (n NodePaths) String() string {
	paths := make([]string, 0, len(n))
	for _, path := range n {
		paths = append(paths, fmt.Sprintf("[%s]", path.String()))
	}
	return strings.Join(paths, ", ")
}

var knownPathCount = make(map[string]int)

func countPaths(s *Node, d *Node) int {
	ps := fmt.Sprintf("%s->%s", s.ID(), d.ID())
	if count, known := knownPathCount[ps]; known {
		return count
	}
	if s.ID() == d.ID() {
		return 1
	}

	count := 0
	for _, o := range s.Nodes() {
		count += countPaths(o, d)
	}

	knownPathCount[ps] = count

	return count
}

func getPaths(s *Node, d *Node) NodePaths {
	var paths NodePaths

	if s.ID() == d.ID() {
		paths = append(paths, NodePath{d})
	} else {
		for _, o := range s.Nodes() {
			for _, oPath := range getPaths(o, d) {
				paths = append(paths, append(NodePath{s}, oPath...))
			}
		}
	}

	return paths
}

type NodePath []*Node

func (np NodePath) String() string {
	ids := make([]string, 0, len(np))
	for _, n := range np {
		ids = append(ids, n.ID())
	}
	return strings.Join(ids, "->")
}

type Nodes map[string]*Node

func (n Nodes) String() string {
	ids := make([]string, 0, len(n))
	for _, n := range n {
		ids = append(ids, fmt.Sprintf("{%s}", n.ID()))
	}
	return strings.Join(ids, ",")
}

func (n Nodes) GetOrCreate(id string) *Node {
	if node, ok := n[id]; ok {
		return node
	}
	n[id] = &Node{id: id, next: make(map[string]*Node)}
	return n[id]
}

type Node struct {
	id   string
	next map[string]*Node
}

func (n *Node) ID() string {
	return n.id
}

func (n *Node) Add(o *Node) {
	n.next[o.id] = o
}

func (n *Node) IsConnectedTo(id string) bool {
	_, connected := n.next[id]
	return connected
}

func (n *Node) Get(id string) *Node {
	if n.IsConnectedTo(id) {
		return n.next[id]
	}
	return nil
}

func (n *Node) Nodes() []*Node {
	var nodes []*Node
	for _, o := range n.next {
		nodes = append(nodes, o)
	}
	return nodes
}

func getNodes(filename string) Nodes {
	replacer := strings.NewReplacer(":", "")
	nodes := make(Nodes)
	for _, line := range files.MustGetLines(filename) {
		ids := strings.Fields(replacer.Replace(strings.TrimSpace(line)))
		n := nodes.GetOrCreate(ids[0])
		for _, id := range ids[1:] {
			o := nodes.GetOrCreate(id)
			n.Add(o)
		}
	}
	return nodes
}
