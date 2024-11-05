package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
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

func start() {
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

		// Read the header
		// buf := make([]byte, calc.HeaderLength)
		// if _, err = client.Read(buf); err != nil {
		// 	log.Panic(err)
		// }

		// Parse header

		dec := gob.NewDecoder(conn)
		var message Message
		dec.Decode(&message)

		fmt.Print("Req: ", message.Request, "\nData: ", string(message.Data))

		if err = conn.Close(); err != nil {
			log.Print(err)
		}
	}
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
