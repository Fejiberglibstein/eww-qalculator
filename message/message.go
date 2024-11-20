package message

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
)

const Port = "/tmp/eww-calc"

type Request uint8

type Message struct {
	Header uint8
	Data   string
}

func SendMessage(conn *net.Conn, message Message) {
	// Allocate a new byte buffer to fill with the bytes of our message
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(message); err != nil {
		log.Panic(err)
	}

	// Write the buffer to the stream
	_, err := (*conn).Write(buf.Bytes())
	if err != nil {
		log.Panic(err)
	}
}

const (
	Listen Request = iota
	SendExpr
)
