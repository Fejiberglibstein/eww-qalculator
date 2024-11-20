package main

import (
	"log"
	"os"

	"github.com/Fejiberglibstein/eww-qalculator/daemon"
	"github.com/Fejiberglibstein/eww-qalculator/listen"
	"github.com/Fejiberglibstein/eww-qalculator/send"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Fatal("Not enough args")
	}

	switch args[1] {
	case "start":
		server, err := daemon.NewServer(args[2:])
		if err != nil {
			log.Panic(err)
		}
		server.Run()
	case "send":
		send.Send(args[2:])
	case "listen":
		listen.Listen(args[2:])
	default:
		log.Print("Invalid arguments")
	}

}
