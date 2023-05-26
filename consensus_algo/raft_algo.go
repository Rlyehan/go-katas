package main
import (
	"fmt"
	"sync"
	"time"
)

type Server struct {
	id 		int
	state 	string
	votes 	int
}

type Cluster struct {
	servers 	[]*Server
	leader 		*Server
	mu 			sync.Mutex
}

func NewServer(id int) *Server {
	return &Server{
		id: id,
		state: "follower",
		votes: 0,
	}
}

func NewCluster(n int) *Cluster {
	cluster := &Cluster{
		servers: make([]*Server, n),
	}

	for i := 0; i < n; i++ {
		cluster.servers[i] = NewServer(i)
	}
	return cluster
}

func (c *Cluster) election() {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Println("Election started...")
	for _, server := range c.servers {
		if server.state == "follower" {
			server.state = "candidate"
		}
		fmt.Printf("Server %d became candidate\n", server.id)
	}

	c.leader = nil
	time.AfterFunc(time.Second, c.assignLeader)
}

func (c *Cluster) assignLeader() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, server := range c.servers {
		if server.state == "candidate" && (c.leader == nil || server.votes > c.leader.votes) {
			c.leader = server
		}
	}

	if c.leader != nil {
		c.leader.state = "leader"
		fmt.Printf("Server %d became leader\n", c.leader.id)
	}
}

func (c *Cluster) run () {
	for {
		if c.leader == nil {
			c.election()
		} else {
			fmt.Printf("Server %d is leader\n", c.leader.id)
			time.Sleep(time.Second)
		}
	}
}

func main() {
	cluster := NewCluster(5)
	go cluster.run()
	time.Sleep(time.Minute)
}
