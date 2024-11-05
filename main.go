package main

import (
	"encoding/gob"
	"log"
	"os"
)

func main() {
	args := os.Args

	switch args[1] {
	case "start":
		start()
		break
	case "send":
		send()
		break
	default:
		log.Print("Invalid arguments")
	}

}

func start() {
	if err := os.Remove(Port); err != nil {
		log.Print(err)
	}

	server, err := net.Listen("unix", Port)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Print("Failed to accept request: ", err)
			continue
		}

		// Read the header
		// buf := make([]byte, calc.HeaderLength)
		// if _, err = client.Read(buf); err != nil {
		// 	log.Panic(err)
		// }

		// Parse header

		dec := gob.NewDecoder(conn)
		var message Message
		dec.Decode(&message)

		fmt.Print("Req: ", message.Request, "\nData: ", string(message.Data))

		if err = conn.Close(); err != nil {
			log.Print(err)
		}
	}
}

func send() {

	if len(args) < 3 {
		log.Fatal("Not enough args")
	}

	conn, err := net.Dial("unix", calc.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Get the message to send based on args
	message, err := get_message(args[1:])
	if err != nil {
		log.Fatal(err)
	}

	// Allocate a new byte buffer to fill with the bytes of our message
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err = enc.Encode(message); err != nil {
		log.Fatal(err)
	}

	// Write the buffer to the stream
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

}

func get_message(args []string) (calc.Message, error) {
	data := []byte(args[1])
	switch args[0] {
	case "expr":
		return calc.Message{Data: data, Request: calc.Expr}, nil
	default:
		return calc.Message{}, errors.New("Invalid request")
	}
}
