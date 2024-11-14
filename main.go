package main

import (
	"log"
	"os"

	"github.com/Fejiberglibstein/eww-qalculator/send"
	"github.com/Fejiberglibstein/eww-qalculator/server"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Fatal("Not enough args")
	}

	switch args[1] {
	case "start":
		server, err := server.NewServer(args[2:])
		if err != nil {
			log.Panic(err)
		}
		server.Run()
	case "send":
		send.Send(args[2:])
	case "listen":
		// TODO
	default:
		log.Print("Invalid arguments")
	}

}
