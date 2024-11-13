package send

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"net"

	"github.com/Fejiberglibstein/eww-qalculator/message"
	"github.com/Fejiberglibstein/eww-qalculator/server"
)

func Send(args []string) {
	if len(args) < 2 {
		log.Fatal("Not enough args for send request")
	}

	conn, err := net.Dial("unix", server.Port)
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

func getMessage(args []string) (message.Message, error) {
	switch args[0] {
	case "expr":
		return message.Message{Data: args[1] + "\n", Header: uint8(message.SendExpr)}, nil
	default:
		return message.Message{}, errors.New("Invalid request")
	}
}
