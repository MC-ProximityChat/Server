package server

import (
	"errors"
	"sync"
)

type Manager struct {
	Servers sync.Map
}

func NewManager() *Manager {
	return &Manager{Servers: sync.Map{}}
}

func (m *Manager) Create(id string, server *Server) {
	m.Servers.Store(id, server)
}

func (m *Manager) Read(id string) (*Server, bool) {
	server, ok := m.Servers.Load(id)
	var serverCast *Server = nil

	if !ok {
		return nil, false
	} else {
		serverCast = server.(*Server)
	}

	return serverCast, true
}

func (m *Manager) Delete(id string) error {

	var err error

	if !m.Contains(id) {
		err = errors.New("server doesn't exist")
	}

	m.Servers.Delete(id)

	return err
}

func (m *Manager) Contains(id string) bool {
	_, ok := m.Servers.Load(id)
	return ok
}
