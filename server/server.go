package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/Fejiberglibstein/eww-qalculator/message"
	"github.com/Fejiberglibstein/eww-qalculator/parser"
)

const Port = "/tmp/eww-calc"

type Server struct {
	listener net.Listener
	qalc     Qalc
}

type Qalc struct {
	cmd        *exec.Cmd
	stdout     *bufio.Reader
	stdin      io.WriteCloser
	stdoutPipe io.ReadCloser
}

func NewServer(args []string) (Server, error) {
	// Print just an empty array to give eww somehting to start with
	fmt.Println("[]")
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

	stdout := bufio.NewReader(stdoutPipe)

	return Qalc{
		cmd:        qalc,
		stdout:     stdout,
		stdin:      stdinPipe,
		stdoutPipe: stdoutPipe,
	}, nil
}

func (s *Server) Run() {
	defer s.listener.Close()

	qalc, err := runQalc()
	if err != nil {
		log.Panic(err)
	}
	s.qalc = qalc

	defer qalc.stdoutPipe.Close()
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
		message, err := message.ReadMessage(&conn)
		if err != nil {
			log.Print(err)
			continue
		}

		if err = s.onRequest(message); err != nil {
			log.Print(err)
		}

	}
}

func (s *Server) onRequest(msg message.Message) error {
	switch msg.Header {
	case uint8(message.SendExpr):
		io.WriteString(s.qalc.stdin, string(msg.Data)+"\n")

		// Read the first line from qalc, this will always be
		//
		// > (whatever expression was inputted)
		//
		// So we can safely ignore it
		if _, err := s.qalc.stdout.ReadString('\n'); err != nil {
			log.Print("Could not read from qalc")
			return err
		}

		// var total string
		qalcStrings := make([]string, 0)
		for {
			// Concatenate all the strings together
			res, err := s.qalc.stdout.ReadString('\n')
			if err != nil {
				log.Print("Could not read from qalc")
				return err
			}
			if res[0] == '>' {
				break
			}

			if res != "\n" {
				qalcStrings = append(qalcStrings, res)
			}
		}

		lines := parser.ParseLines(qalcStrings)
		str, err := json.Marshal(&lines)
		if err != nil {
			log.Print("Could not parse json for this", err)
		}
		fmt.Println(string(str))

	default:
		return errors.New("Invalid request received")
	}
	return nil
}
