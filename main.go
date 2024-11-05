package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"go/scanner"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Fatal("Not enough args")
	}

	switch args[1] {
	case "start":
		start()
		break
	case "send":
		send(args[2:])
		break
	default:
		log.Print("Invalid arguments")
	}

}

type CalcProcess struct {
}

func start() {
	qalc := exec.Command("qalc", "-terse")

	stdin, err := qalc.StdinPipe()
	if err != nil {
		log.Panic(err)
	}
	defer stdin.Close()

	stdout, err := qalc.StdoutPipe()
	if err != nil {
		log.Panic(err)
	}
	defer stdout.Close()
	reader := bufio.NewReader(stdout)

	if err = qalc.Start(); err != nil {
		log.Panic(err)
	}

	runServer(func(message Message) error {
		switch message.Request {
		case Expr:
			io.WriteString(stdin, string(message.Data))
			// Read the first line from qalc, this is always an empty line
			if _, err = reader.ReadString('\n'); err != nil {
				log.Print("Could not read from qalc")
				return err
			}
			var res string
			for res != "\n" {
				res, err = reader.ReadString('\n')
				if err != nil {
					log.Print("Could not read from qalc")
					return err
				}
			}
		}
		return nil
	})
}

func send(args []string) {
	if len(args) < 2 {
		log.Fatal("Not enough args for send request")
	}

	conn, err := net.Dial("unix", Port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Get the message to send based on args
	message, err := getMessage(args)
	if err != nil {
		log.Fatal(err)
	}

	// Allocate a new byte buffer to fill with the bytes of our message
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err = enc.Encode(message); err != nil {
		log.Fatal(err)
	}

	// Write the buffer to the stream
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

}

func getMessage(args []string) (Message, error) {
	data := []byte(args[1])
	switch args[0] {
	case "expr":
		return Message{Data: data, Request: Expr}, nil
	default:
		return Message{}, errors.New("Invalid request")
	}
}

func runServer(onRequest func(Message) error) {
	if err := os.Remove(Port); err != nil {
		log.Print(err)
	}

	server, err := net.Listen("unix", Port)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Print("Failed to accept request: ", err)
			continue
		}

		// Parse the data
		dec := gob.NewDecoder(conn)
		var message Message
		dec.Decode(&message)

		if err = onRequest(message); err != nil {
			log.Print(err)
		}

		if err = conn.Close(); err != nil {
			log.Print(err)
		}
	}
}
