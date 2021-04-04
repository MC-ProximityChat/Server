package server

import (
	"sync"
	"time"
)

type Manager struct {
	Servers sync.Map
}

type managerJanitor struct {
	manager *Manager
	timer   *time.Timer
}

func NewManager() *Manager {
	return &Manager{Servers: sync.Map{}}
}

func (m *Manager) Add(id string, server *Server) {
	m.Servers.Store(id, server)
}

func (m *Manager) Contains(id string) bool {
	_, ok := m.Servers.Load(id)
	return ok
}

func (m *Manager) Get(id string) (*Server, bool) {
	server, ok := m.Servers.Load(id)
	var serverCast *Server = nil

	if !ok {
		return nil, false
	} else {
		serverCast = server.(*Server)
	}

	return serverCast, true
}

func StartCleanup() {
	go func() {

	}()
}
