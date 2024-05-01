package consistent

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// Node represents a node in the ring
type Node struct {
	Id     string
	HashId uint32
}

// NewNode creates a new Node and returns a pointer to it
func NewNode(id string) *Node {
	return &Node{
		Id:     id,
		HashId: hashId(id),
	}
}

// Nodes is a slice of pointers to Node
type Nodes []*Node

// Len returns the number of nodes in the slice
func (n Nodes) Len() int {
	return len(n)
}

// Ring represents a consistent hash ring
type Ring struct {
	Nodes Nodes
	sync.RWMutex
}

// NewRing creates a new Ring and returns a pointer to it
func NewRing() *Ring {
	return &Ring{Nodes: Nodes{}}
}

// AddNode adds a new node to the ring
func (r *Ring) AddNode(id string) {
	r.Lock()
	defer r.Unlock()

	node := NewNode(id)
	r.Nodes = append(r.Nodes, node)

	sort.SliceStable(r.Nodes, func(i, j int) bool {
		return r.Nodes[i].HashId < r.Nodes[j].HashId
	})
	// // 	 Print all the nodes in the ring
	// nodes := getNodesList(r.Nodes)
	// println("Nodes in the ring: ", nodes)
}

// Get flattened list of nodes in the ring
func getNodesList(r Nodes) string {
	var nodes string
	for _, node := range r {
		nodes = nodes + strconv.Itoa(int(node.HashId)) + ","
	}
	return nodes
}

// RemoveNode removes a node from the ring by its id
func (r *Ring) RemoveNode(id string) error {
	r.Lock()
	defer r.Unlock()

	nodeIndex := r.search(id)
	if nodeIndex >= r.Nodes.Len() || r.Nodes[nodeIndex].Id != id {
		return errors.New("node not found")
	}

	r.Nodes = append(r.Nodes[:nodeIndex], r.Nodes[nodeIndex+1:]...)

	return nil
}

// Get returns the id of the node that should handle the request for the given id
func (r *Ring) Get(id string) (string, string, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.Nodes) == 0 {
		return "", "", errors.New("no nodes available")
	}

	nodeIndex := r.search(id)
	if nodeIndex >= r.Nodes.Len() {
		nodeIndex = 0
	}

	// If the node is not available, return the next node
	if r.Nodes[nodeIndex].Id != id {
		newNodeIndex := (nodeIndex + 1) % r.Nodes.Len()
		return r.Nodes[nodeIndex].Id, r.Nodes[newNodeIndex].Id, nil
	}

	return r.Nodes[nodeIndex].Id, "", nil
}

// search returns the index of the node that should handle the request for the given id
func (r *Ring) search(id string) int {
	return sort.Search(r.Nodes.Len(), func(i int) bool {
		return r.Nodes[i].HashId >= hashId(id)
	})
}

// hashId generates a hash for the given key
func hashId(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}
