package server

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Manager struct {
	Servers sync.Map
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

func (m *Manager) Get(id string) *Server {
	server, ok := m.Servers.Load(id)
	var serverCast *Server = nil

	if !ok {
		logrus.Errorf("Couldn't find server %s", id)
	} else {
		serverCast = server.(*Server)
	}

	return serverCast
}
