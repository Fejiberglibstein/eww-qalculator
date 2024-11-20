package listen

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			msg, err := message.ReadMessage(&conn)
			if err != nil {
				if err == io.EOF {
					c <- nil
				}
				continue
			}
			fmt.Println(msg.Data)
		}
	}()
	<-c

	log.Print(conn.Close())
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
