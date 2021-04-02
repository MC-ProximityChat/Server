package server

import (
	routing "github.com/jackwhelpton/fasthttp-routing"
	"github.com/sirupsen/logrus"
)

type Server struct {
	DataQueue chan LocationPacket
	CloseChan chan bool
}

func NewServer() *Server {
	return &Server{
		DataQueue: make(chan LocationPacket),
		CloseChan: make(chan bool),
	}
}

func (s *Server) Close() {
	s.CloseChan <- true
}

func (s *Server) SendPacket(packet LocationPacket) {
	s.DataQueue <- packet
}

func (s *Server) Run() {
	select {
	case <-s.DataQueue:
		return
	case <-s.CloseChan:
		return
	}
}

type LocationPacket struct {
	UUID string  `json:"uuid"`
	X    float64 `json:"x"`
	Z    float64 `json:"z"`
}

func NewLocationPacket(context *routing.Context) *LocationPacket {
	var locationPacket LocationPacket

	if err := context.Read(&locationPacket); err != nil {
		logrus.Fatalf("Unable to read JSON %s", err)
	}

	return &locationPacket
}
