package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/bbengfort/whisper"
)

func main() {

	listen, err := net.Listen("tcp", "localhost:3264")
	if err != nil {
		log.Fatal(err)
	}

	defer listen.Close()

	for {
		c, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go serve(c)
	}
}

func serve(c net.Conn) {
	defer c.Close()

	dec := json.NewDecoder(c)
	var m whisper.Message
	if err := dec.Decode(&m); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stdout, "%s: \"%s\"\n", m.Sender, m.Body)
	fmt.Fprintf(c, "Received message at %s\n", m.Timestamp)
	io.Copy(c, c)

}
