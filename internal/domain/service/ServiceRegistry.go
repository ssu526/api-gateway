package service

type ServiceRegistry interface {
	RegisterService(service *Service) error
	DeregisterService(name string) error
	GetService(name string) (*Service, error)
}
