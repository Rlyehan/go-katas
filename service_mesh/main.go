package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	// Create the service registry and register the services
	registry := &ServiceRegistry{services: make(map[string]Service)}
	registry.Register("serviceA", &ServiceA{})
	registry.Register("serviceB", &ServiceB{})

	// Create the service proxy
	proxy := &ServiceProxy{registry: registry}

	// Run the service proxy in its own goroutine, so it doesn't block the rest of the program
	go func() {
		log.Println("Service proxy is running on port 8080")
		if err := http.ListenAndServe(":8080", proxy); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait a moment to ensure the server has time to start up
	time.Sleep(1 * time.Second)

	// Make requests to the services
	testService("serviceA")
	testService("serviceB")
}

func testService(serviceName string) {
	resp, err := http.Get("http://localhost:8080/" + serviceName + "?number=5")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s response: %s", serviceName, body)
}
