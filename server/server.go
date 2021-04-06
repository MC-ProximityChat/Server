package server

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/sony/sonyflake"
	"proximity-chat-grouper/util"
	"strconv"
	"time"
)

var (
	serverIdGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})
	expirationPolicy  = 2 * time.Minute
)

// Server
type Server struct {
	ID                  string
	Name                string
	PotentialUsersCache *util.BiCache
	Users               []string
	DataQueue           chan Location
	CloseChan           chan bool
}

type PotentialUser struct {
	Code             string `json:"code"`
	ExpirationPolicy string `json:"expirationPolicy"`
	AlreadyGenerated bool   `json:"alreadyGenerated"`
}

func NewPotentialUser() (*PotentialUser, error) {

	code, err := util.GenerateRandomCode()

	if err != nil {
		return nil, err
	}

	return &PotentialUser{
		Code:             code,
		ExpirationPolicy: expirationPolicy.String(),
		AlreadyGenerated: false,
	}, nil
}

// Creates new server object
func NewServer(name string) *Server {
	id, err := serverIdGenerator.NextID()
	if err != nil {
		logrus.Errorf("Unable to generate uuid %s", err)
	}

	idStr := strconv.FormatUint(id, 10)

	return &Server{
		ID:                  idStr,
		Name:                name,
		PotentialUsersCache: util.NewBiCache(expirationPolicy, 10*time.Minute),
		Users:               make([]string, 0, 10),
		DataQueue:           make(chan Location),
		CloseChan:           make(chan bool),
	}
}

func (s *Server) CreatePotentialUser(uuid string) (*PotentialUser, error) {

	if s.PotentialUsersCache.ContainsValue(uuid) { // user already generated a cache
		logrus.Info("here")
		code, expirationTime, ok := s.PotentialUsersCache.GetKeyWithExpiration(uuid)
		duration := expirationTime.Sub(time.Now()).Round(1 * time.Second)
		if !ok {
			return nil, errors.New("internal server error (ok somethings really fucked up lmao)")
		}

		return &PotentialUser{
			Code:             code,
			ExpirationPolicy: duration.String(),
			AlreadyGenerated: true,
		}, nil
	}
	potentialUser, err := NewPotentialUser()

	logrus.Info(uuid)

	if err != nil {
		return nil, err
	}

	s.PotentialUsersCache.Set(potentialUser.Code, uuid, cache.DefaultExpiration)

	return potentialUser, nil
}

func (s *Server) getUuidFromCode(code string) (string, bool) {
	return s.PotentialUsersCache.GetValueAndDelete(code)
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
