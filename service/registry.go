package service

func NewRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services:   make([]*ButlerService, 0, 5),
		serviceMap: make(map[string]*ButlerService),
	}
}

type ServiceRegistry struct {
	services   []*ButlerService
	serviceMap map[string]*ButlerService
}

func (r *ServiceRegistry) AddService(service *ButlerService) {
	r.services = append(r.services, service)
	r.serviceMap[service.Name] = service
}

func (r *ServiceRegistry) List() []*ButlerService {
	return r.services
}

func (r *ServiceRegistry) GetByName(name string) (*ButlerService, bool) {
	service, ok := r.serviceMap[name]
	return service, ok
}
