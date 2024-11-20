package message

import (
	"bytes"
	"encoding/gob"
	"net"
)

const Port = "/tmp/eww-calc"

type Request uint8

type Message struct {
	Header uint8
	Data   string
}

func ReadMessage(conn *net.Conn) (Message, error) {
	dec := gob.NewDecoder(*conn)
	var message Message
	if err := dec.Decode(&message); err != nil {
		return Message{}, err
	}

	return message, nil
}

func SendMessage(conn *net.Conn, message Message) error {
	// Allocate a new byte buffer to fill with the bytes of our message
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(message); err != nil {
		return err
	}

	// Write the buffer to the stream
	_, err := (*conn).Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

const (
	// Data is expected to be the channel name you want to listen on
	Listen Request = iota
	// Data is expected to be the expression to send
	SendExpr
)
