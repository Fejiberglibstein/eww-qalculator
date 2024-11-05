package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"net"
	"os"

	calc "github.com/Fejiberglibstein/eww-qalculator"
)

func main() {
	args := os.Args

	if len(args) < 3 {
		log.Fatal("Not enough args")
	}

	conn, err := net.Dial("unix", calc.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Get the message to send based on args
	message, err := get_message(args[1:])
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

func get_message(args []string) (calc.Message, error) {

	data := []byte(args[1])
	switch args[0] {
	case "expr":
		return calc.Message{Data: data, Request: calc.Expr}, nil
	default:
		return calc.Message{}, errors.New("Invalid request")
	}
}
