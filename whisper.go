// Package whisper provides the library functionality for the P2P whispernet
// application written in response to the tutorial by Andrew Gerrand and
// Francesc Campoy at http://whispering-gophers.appspot.com/talk.slide.
package whisper

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

// Version specifies the current revision of the whisper library
const Version = "1.0"

// Message represents data sent as JSON across the whispernet
type Message struct {
	Sender string    `json:"sender"` // The name of the message sender
	Addr   string    `json:"addr"`   // The address of the sender
	Body   string    `json:"body"`   // The body of the message
	SentAt time.Time `json:"sent"`   // The time that the message was sent
}

// Print a message to a representative string.
func (msg Message) String() string {
	return fmt.Sprintf("[%s] %s: %s", msg.SentAt.Format(time.Stamp), msg.Sender, msg.Body)
}

// Client is a whisper agent that accepts user input and sends messages.
type Client struct {
	Name    string        // The user name of the client
	Input   *InputHandler // The handler for user input from the console
	Address string        // Address to listen on for messages.
	Server  string        // Address of the server to send data to

}

// NewClient constructs a client and instantiates handlers.
func NewClient(name string, server string) *Client {
	return &Client{
		Name:   name,
		Input:  NewInputHandler(">"),
		Server: server,
	}
}

// Connect to the given server address
func (client *Client) Connect() (net.Conn, *Error) {
	conn, err := net.Dial("tcp", client.Server)
	if err != nil {
		return nil, &Error{fmt.Sprintf("Could not connect to %s: %s", client.Server, err.Error()), 99}
	}

	return conn, nil
}

// Run the handler and sends any messages from the command line.
func (client *Client) Run() *Error {

	err := make(chan *Error)

	// Run the listener handler
	go client.Listen(err)

	// Now handle all user input
	go client.Handle(err)

	return <-err
}

// Listen accepts incomming connections and prints messages to the console.
func (client *Client) Listen(echan chan<- *Error) {
	addr, iperr := ExternalIP()
	if iperr != nil {
		echan <- iperr
		close(echan)
		return
	}

	listen, err := net.Listen("tcp4", addr+":0")
	if err != nil {
		echan <- &Error{fmt.Sprintf("Couldn't listen on %s:0 -- %s", addr, err.Error()), 99}
		close(echan)
		return
	}

	defer listen.Close()
	client.Address = listen.Addr().String()

	// Mention what is going on, on which address.
	fmt.Printf("Listening on %s, quit by typing exit\n> ", client.Address)

	for {
		conn, err := listen.Accept()
		if err != nil {
			echan <- &Error{fmt.Sprintf("Couldn't accept connection: %s", err.Error()), 98}
			close(echan)
			return
		}

		go client.Recv(conn)
	}
}

// Send constructs a message object for JSON serialization and puts it on the
// TCP connection object that the client maintains as open with the server.
func (client *Client) Send(body string) *Error {

	msg := &Message{
		Sender: client.Name,
		Addr:   client.Address,
		Body:   body,
		SentAt: time.Now(),
	}

	conn, err := client.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	enc := json.NewEncoder(conn)
	if err := enc.Encode(msg); err != nil {
		return &Error{fmt.Sprintf("Could not encode message: %s", err.Error()), 3}
	}

	fmt.Fprintf(os.Stdout, "\r\r> %s\n", msg.String())

	return nil
}

// Recv deserializes JSON messages from the stream and prints them out.
func (client *Client) Recv(conn net.Conn) *Error {
	defer conn.Close()

	dec := json.NewDecoder(conn)
	var m Message
	if err := dec.Decode(&m); err != nil {
		return &Error{fmt.Sprintf("Could not decode message: %s", err.Error()), 4}
	}

	fmt.Fprintf(os.Stdout, "> Message from %s", m.Addr)
	fmt.Fprintf(os.Stdout, "\r> %s\n> ", m.String())
	return nil
}
