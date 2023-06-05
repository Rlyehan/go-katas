package main

type ServiceRegistry struct {
	services map[string]Service
}

func (s *ServiceRegistry) Register(name string, service Service) {
	s.services[name] = service
}

func (s *ServiceRegistry) Get(name string) Service {
	return s.services[name]
}
