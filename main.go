package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/Fejiberglibstein/eww-qalculator/message"
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

func start() {
	qalc := exec.Command("qalc")

	stdinPipe, err := qalc.StdinPipe()
	if err != nil {
		log.Panic(err)
	}
	defer stdinPipe.Close()

	stdoutPipe, err := qalc.StdoutPipe()
	if err != nil {
		log.Panic(err)
	}
	stdout := bufio.NewReader(stdoutPipe)
	defer stdoutPipe.Close()

	if err = qalc.Start(); err != nil {
		log.Panic(err)
	}
	defer qalc.Process.Kill()

}

func runServer(onRequest func(Message) error) {
}
