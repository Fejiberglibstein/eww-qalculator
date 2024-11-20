package listen

import (
	"fmt"
	"log"
	"net"

	"github.com/Fejiberglibstein/eww-qalculator/message"
)

type Channel string

const (
	ExprChan   Channel = "expr"
	ResultChan Channel = "result"
)

type listener struct {
	channel Channel
	conn    net.Conn
}

func Listen(args []string) {
	channel := Channel(args[0])

	conn, err := net.Dial("unix", message.Port)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	l := listener{
		channel: channel,
		conn:    conn,
	}

	l.sendListenReq(channel)

	for {
		msg, err := message.ReadMessage(&conn)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			continue
		}
		fmt.Print(msg.Data)
	}
}

func (l *listener) sendListenReq(channel Channel) {
	message.SendMessage(
		&l.conn,
		message.Message{
			Header: uint8(message.Listen),
			Data:   string(channel),
		},
	)
}
