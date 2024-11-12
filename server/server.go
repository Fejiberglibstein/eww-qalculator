package server

import (
	"encoding/gob"
	"errors"
	"log"
	"net"
	"os"

	"github.com/Fejiberglibstein/eww-qalculator/message"
)

const Port = "/tmp/eww-calc"

type Server struct {
	listener net.Listener
}

func NewServer() (Server, error) {
	if err := os.Remove(Port); err != nil {
		log.Print(err)
	}

	listener, err := net.Listen("unix", Port)
	if err != nil {
		return Server{}, err
	}

	return Server{
		listener: listener,
	}, nil

}

func (s *Server) Run() {
	defer s.listener.Close()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Print("Failed to accept request: ", err)
			continue
		}

		// Parse the data
		dec := gob.NewDecoder(conn)
		var message message.Message
		dec.Decode(&message)

		if err = onRequest(message); err != nil {
			log.Print(err)
		}

		if err = conn.Close(); err != nil {
			log.Print(err)
		}
	}
}
