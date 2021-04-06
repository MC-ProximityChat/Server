package server

import "errors"

type Service struct {
	manager *Manager
}

type SimplifiedServer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewService() *Service {
	return &Service{
		manager: NewManager(),
	}
}

func (s *Service) CreateServer(name string) *SimplifiedServer {

	server := NewServer(name)
	s.manager.Create(server.ID, server)

	return &SimplifiedServer{
		ID:   server.ID,
		Name: server.Name,
	}
}

func (s *Service) ReadServer(id string) (*Server, error) {
	server, ok := s.manager.Read(id)

	var err error

	if !ok {
		err = errors.New("server not found")
	}

	return server, err
}

func (s *Service) ReadServerAsSimplified(id string) (*SimplifiedServer, error) {
	server, ok := s.manager.Read(id)

	var err error
	var simplifiedServer *SimplifiedServer

	if !ok {
		err = errors.New("server not found")
	} else {
		simplifiedServer = &SimplifiedServer{ID: server.ID, Name: server.Name}
	}

	return simplifiedServer, err
}

func (s *Service) DeleteServer(id string) error {
	return s.manager.Delete(id)
}
