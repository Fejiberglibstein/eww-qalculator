package listen

import (
	"fmt"
	"log"
	"net"

	"github.com/Fejiberglibstein/eww-qalculator/message"
	"github.com/Fejiberglibstein/eww-qalculator/server"
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

func Listen(channel Channel) {
	conn, err := net.Dial("unix", server.Port)
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
			log.Print(err)
			continue
		}

		fmt.Print(msg)
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
