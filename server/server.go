package server

import (
	routing "github.com/jackwhelpton/fasthttp-routing"
	"github.com/sirupsen/logrus"
	"github.com/sony/sonyflake"
	"strconv"
)

var (
	idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})
)

// Server
type Server struct {
	ID        string
	Name      string
	DataQueue chan Location
	CloseChan chan bool
}

// Creates new server object
func NewServer(name string) *Server {
	id, err := idGenerator.NextID()
	if err != nil {
		logrus.Errorf("Unable to generate uuid %s", err)
	}

	idStr := strconv.FormatUint(id, 10)

	return &Server{
		ID:        idStr,
		Name:      name,
		DataQueue: make(chan Location),
		CloseChan: make(chan bool),
	}
}

// Closes channel
func (s *Server) Close() {
	s.CloseChan <- true
}

// Sends location to channel
func (s *Server) SendLocation(packet Location) {
	s.DataQueue <- packet
}

// Listens to channels
func (s *Server) Run() {
	select {
	case <-s.DataQueue:
		return
	case <-s.CloseChan:
		return
	}
}

// Location information sent from client
type Location struct {
	UUID string  `json:"uuid"`
	X    float64 `json:"x"`
	Z    float64 `json:"z"`
}

// Creates new location from routing context (ie. POST body)
func NewLocation(context *routing.Context) *Location {
	var locationPacket Location

	if err := context.Read(&locationPacket); err != nil {
		logrus.Fatalf("Unable to read JSON %s", err)
	}

	return &locationPacket
}
