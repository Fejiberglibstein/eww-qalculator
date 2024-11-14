package server

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/Fejiberglibstein/eww-qalculator/message"
)

const Port = "/tmp/eww-calc"

type Server struct {
	listener net.Listener
}

type Qalc struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	stdin  io.WriteCloser
}

func NewServer(args []string) (Server, error) {
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

func runQalc() (Qalc, error) {
	qalc := exec.Command("qalc")

	stdinPipe, err := qalc.StdinPipe()
	if err != nil {
		return Qalc{}, err
	}

	stdoutPipe, err := qalc.StdoutPipe()
	if err != nil {
		return Qalc{}, err
	}

	if err = qalc.Start(); err != nil {
		return Qalc{}, err
	}

	return Qalc{
		cmd:    qalc,
		stdout: stdoutPipe,
		stdin:  stdinPipe,
	}, nil
}

func (s *Server) Run() {
	defer s.listener.Close()

	qalc, err := runQalc()
	if err != nil {
		log.Panic(err)
	}
	defer qalc.stdout.Close()
	defer qalc.stdin.Close()

	s.listen()
}

func (s *Server) listen() {
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

		if err = s.onRequest(message); err != nil {
			log.Print(err)
		}

		if err = conn.Close(); err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) onRequest(msg message.Message) error {
	switch msg.Header {
	case uint8(message.SendExpr):
		io.WriteString(stdinPipe, string(msg.Data)+"\n")

		// Read the first line from qalc, this will always be
		//
		// > (whatever expression was inputted)
		//
		// So we can safely ignore it
		if _, err = stdout.ReadString('\n'); err != nil {
			log.Print("Could not read from qalc")
			return err
		}

		// var total string
		var res string = " "
		var total string
		for {
			// Concatenate all the strings together
			res, err = stdout.ReadString('\n')
			if err != nil {
				log.Print("Could not read from qalc")
				return err
			}
			if res[0] == '>' {
				break
			}
			total += res
		}
		fmt.Print(total)
	default:
		return errors.New("Invalid request received")
	}
	return nil
	return nil
}
