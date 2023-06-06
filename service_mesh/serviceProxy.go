package main

import (
	"log"
	"net/http"
	"strings"
)

type ServiceProxy struct {
	registry *ServiceRegistry
}

func (p *ServiceProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		http.Error(w, "invalid service request", http.StatusBadRequest)
		return
	}

	serviceName := parts[1]
	service := p.registry.Get(serviceName)
	if service == nil {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}
	log.Println("Proxying request to", serviceName)

	r.URL.Path = "/" + strings.Join(parts[2:], "/")

	service.ServeHTTP(w, r)
}
