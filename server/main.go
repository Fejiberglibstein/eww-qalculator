package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Fejiberglibstein/eww-qalculator"
)

func main() {
	if err := os.Remove(calc.Port); err != nil {
		log.Print(err)
	}

	server, err := net.Listen("unix", calc.Port)
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
		var message calc.Message
		dec.Decode(&message)

		fmt.Print("Req: ", message.Request, "\nData: ", string(message.Data))

		if err = conn.Close(); err != nil {
			log.Print(err)
		}
	}
}
