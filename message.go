package main

const Port = "/tmp/eww-calc"

type Request uint8

type Message struct {
	Request Request
	Data    []byte
}

const (
	Expr Request = iota
)
