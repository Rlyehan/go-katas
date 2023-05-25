package main

import ( 
	"hash/fnv"
	"sort"
	"sync"
	"strconv"
)

type Cache struct {
	m 	map[string]string
	mux sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		m: make(map[string]string),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	value, ok := c.m[key]
	return value, ok
}

func (c *Cache) Set(key, value string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.m[key] = value
}

func (c *Cache) Delete(key string) {
	 c.mux.Lock()
	 defer c.mux.Unlock()
	 delete(c.m, key)
}

func Hash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32())
}

type Node struct {
	Cache 	*Cache
	Hash	int
}

type Ring struct {
	Nodes []*Node
}

func NewRing(nodeNames []string) *Ring {
	nodes := make([]*Node, len(nodeNames))
	for i, name := range nodeNames {
		nodes[i] = &Node{
			Cache: NewCache(),
			Hash:  Hash(name),
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Hash < nodes[j].Hash
	})
	return &Ring{Nodes: nodes}
}

func (r *Ring) GetNode(key string) *Node {
	keyHash := Hash(key)
	for _, node := range r.Nodes {
		if node.Hash >= keyHash {
			return node
		}
	}
	return r.Nodes[0]
}

func (r *Ring) Get(key string) (string, bool) {
	node := r.GetNode(key)
	return node.Cache.Get(key)
}

func (r *Ring) Set(key, value string) {
	node := r.GetNode(key)
	node.Cache.Set(key, value)
}

func (r *Ring) Delete(key string) {
	node := r.GetNode(key)
	node.Cache.Delete(key)
}

func (r *Ring) GetNodeName(key string) string {
	node := r.GetNode(key)
	for i, n := range r.Nodes {
		if n == node {
			return "node" + strconv.Itoa(i+1)
		}
	}
	return ""
}

func main() {
	nodeNames := []string{"alpha", "beta", "gamma"}
	ring := NewRing(nodeNames)

	ring.Set("hello", "world")
	ring.Set("foo", "bar")
	ring.Set("fizz", "buzz")
	ring.Set("general", "kenobi")
	ring.Set("infinite", "loop")

	keys := []string{"hello", "foo", "fizz", "general", "infinite"}
	for _, key := range keys {
		println("Key:", key, "is on", ring.GetNodeName(key))
	}

	value, _ := ring.Get("hello")
	println(value)

	value, _ = ring.Get("foo")
	println(value)

	value, _ = ring.Get("fizz")
	println(value)

	value, _ = ring.Get("general")
	println(value) 

	value, _ = ring.Get("infinite")
	println(value) 

	ring.Delete("general")

	value, found := ring.Get("general")
	if !found {
		println("general not found")
	}
}

