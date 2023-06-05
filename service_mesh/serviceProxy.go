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
	// Extract service name from the request, e.g., /serviceA/
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		http.Error(w, "invalid service request", http.StatusBadRequest)
		return
	}

	// Get the service from the registry
	serviceName := parts[1]
	service := p.registry.Get(serviceName)
	if service == nil {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}
	log.Println("Proxying request to", serviceName)

	// Update the request to remove the service name
	r.URL.Path = "/" + strings.Join(parts[2:], "/")

	// Forward the request to the service
	service.ServeHTTP(w, r)
}
