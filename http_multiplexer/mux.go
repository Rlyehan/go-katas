package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var sem = make(chan struct{}, 20)
var wg sync.WaitGroup
var handlerWg sync.WaitGroup

func handler(w http.ResponseWriter, r *http.Request) {
	sem <- struct{}{}
	defer func() { <-sem }()

	log.Printf("Processing request for %s\n", r.URL.Path)

	fmt.Fprintf(w, "Hello, %q", r.URL.Path)

	log.Printf("Finished processing request for %s\n", r.URL.Path)

	wg.Done()
	handlerWg.Done()
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/foo", handler)
	http.HandleFunc("/bar", handler)

	go func() {
		log.Println("Starting server at :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	time.Sleep(1 * time.Second)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		handlerWg.Add(1)
		go func(i int) {
			log.Printf("Sending request %d\n", i)
			resp, err := http.Get("http://localhost:8080/")
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			log.Printf("Finished request %d\n", i)
		}(i)
	}

	wg.Wait()
	log.Println("All requests sent")
	
	handlerWg.Wait()
	close(sem)

	log.Println("All requests processed")
}
