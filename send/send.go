package send

import (
	"errors"
	"log"
	"net"

	"github.com/Fejiberglibstein/eww-qalculator/message"
)

func Send(args []string) {
	if len(args) < 2 {
		log.Panic("Not enough args for send request")
	}

	conn, err := net.Dial("unix", message.Port)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	// Get the msg to send based on args
	msg, err := getMessage(args)
	if err != nil {
		log.Panic(err)
	}

	if err = message.SendMessage(&conn, msg); err != nil {
		log.Panic(err)
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
