package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type ServiceA struct{}

func (s *ServiceA) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Service A received request")
	params := r.URL.Query()
	numberStr := params.Get("number")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "%d", number+1)
}

type ServiceB struct{}

func (s *ServiceB) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Service B received request")
	params := r.URL.Query()
	numberStr := params.Get("number")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "%d", number+1)
}

