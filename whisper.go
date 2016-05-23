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
	Sender    string    `json:"sender"`    // The name of the message sender
	Body      string    `json:"body"`      // The body of the message
	Timestamp time.Time `json:"timestamp"` // The time that the message was sent
}

// Client is a whisper agent that accepts user input and sends messages.
type Client struct {
	Name  string        // The user name of the client
	Input *InputHandler // The handler for user input from the console
	Conn  net.Conn      // Connection to send the data to.
	// Conns []*net.Conn   // Connections handled by the client

}

// NewClient constructs a client and instantiates handlers.
func NewClient(name string) *Client {
	return &Client{
		Name:  name,
		Input: NewInputHandler(">"),
	}
}

// Connect to the given server address
func (client *Client) Connect(address string) *Error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return &Error{fmt.Sprintf("Could not connect to %s: %s", address, err.Error()), 99}
	}

	// client.Conns = append(client.Conns, &conn)
	client.Conn = conn
	return nil
}

// Close the connection(s) that are maintained by the client
func (client *Client) Close() *Error {
	err := client.Conn.Close()
	if err != nil {
		return &Error{fmt.Sprintf("Could not close connection: %s", err.Error()), 98}
	}
	return nil
}

// Run the handler and sends any messages from the command line.
func (client *Client) Run() *Error {
	for {
		body, err := client.Input.ReadLine()
		if err != nil {
			return err
		}

		// Send the message to the server.
		err = client.Send(body)
		if err != nil {
			return err
		}
	}
}

// Send constructs a message object for JSON serialization and puts it on the
// TCP connection object that the client maintains as open with the server.
func (client *Client) Send(body string) *Error {
	msg := &Message{
		Sender:    client.Name,
		Body:      body,
		Timestamp: time.Now(),
	}

	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(msg); err != nil {
		return &Error{fmt.Sprintf("Could not encode message: %s", err.Error()), 3}
	}

	return nil
}
