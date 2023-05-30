package main

import (
	"fmt"
	"strings"
	"sync"
)

type KeyValue struct {
	Key   string
	Value int
}

func Map(filename string, content string, output chan<- KeyValue) {
	words := strings.Fields(content)
	for _, word := range words {
		output <- KeyValue{word, 1}
	}
	fmt.Println("Map Finished:", filename)
}

func Reduce(key string, values <-chan int, output chan<- KeyValue, wg *sync.WaitGroup) {
	defer wg.Done()
	count := 0
	for v := range values {
		count += v
	}
	output <- KeyValue{key, count}
	fmt.Println("Reduce Finished:", key)
}

func main() {
	docs := map[string]string{
		"file1.txt": "hello world hello",
		"file2.txt": "world",
		"file3.txt": "hello world",
	}

	mapChan := make(chan KeyValue)
	reduceChan := make(chan KeyValue)
	reduceInputChans := make(map[string]chan int)

	var mapWg, reduceWg sync.WaitGroup
	mapWg.Add(len(docs))

	for filename, content := range docs {
		go func(filename string, content string) {
			Map(filename, content, mapChan)
			mapWg.Done()
		}(filename, content)
	}

	go func() {
		mapWg.Wait()
		close(mapChan)
		fmt.Println("All Map jobs are done")
	}()

	go func() {
		for kv := range mapChan {
			if _, ok := reduceInputChans[kv.Key]; !ok {
				ch := make(chan int)
				reduceInputChans[kv.Key] = ch
				reduceWg.Add(1)
				go Reduce(kv.Key, ch, reduceChan, &reduceWg)
			}
			reduceInputChans[kv.Key] <- kv.Value
		}

		for _, ch := range reduceInputChans {
			close(ch)
		}
		fmt.Println("All reduce input channels are closed")
	}()

	go func() {
		reduceWg.Wait()
		close(reduceChan)
		fmt.Println("All Reduce jobs are done")
	}()

	var printWg sync.WaitGroup
	printWg.Add(1)
	go func() {
		defer printWg.Done()
		for kv := range reduceChan {
			fmt.Printf("Key: %s, Value: %d\n", kv.Key, kv.Value)
		}
	}()

	printWg.Wait()

	fmt.Println("End of main function")
}
