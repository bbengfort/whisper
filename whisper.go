// Package whisper provides the library functionality for the P2P whispernet
// application written in response to the tutorial by Andrew Gerrand and
// Francesc Campoy at http://whispering-gophers.appspot.com/talk.slide.
package whisper

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Version specifies the current revision of the whisper library
const Version = "1.0"

// Strings to specify exit queries on input.
const (
	EXIT = "exit"
	QUIT = "quit"
)

// Message represents data sent as JSON across the whispernet
type Message struct {
	Sender    string    `json:"sender"`    // The name of the message sender
	Body      string    `json:"body"`      // The body of the message
	Timestamp time.Time `json:"timestamp"` // The time that the message was sent
}

// NewMessage constructs a message object for JSON serialization
func NewMessage(body string, sender string) *Message {
	return &Message{
		Sender:    sender,
		Body:      body,
		Timestamp: time.Now(),
	}
}

// InputHandler provides a method for reading information from standard input
// with a prompt and dealing with it downstream by calling next on the buffer.
type InputHandler struct {
	prompt string         // The prompt symbol to use
	reader *bufio.Scanner // The scanner to read standard input from
	exit   map[string]int // Exit codes and queries
}

// NewInputHandler creates an input handler connected to standard input.
func NewInputHandler(prompt string) *InputHandler {
	exit := make(map[string]int)
	exit[EXIT] = 1
	exit[QUIT] = 2

	return &InputHandler{
		prompt: prompt,
		reader: bufio.NewScanner(os.Stdin),
		exit:   exit,
	}
}

// ReadLine returns the next line from the buffered reading of stdin.
func (handler *InputHandler) ReadLine() (string, *Error) {
	fmt.Print(handler.prompt + " ")

	if handler.reader.Scan() {
		// Read text and strip stpaces
		output := handler.reader.Text()
		output = strings.TrimSpace(output)

		if _, contains := handler.exit[strings.ToLower(output)]; contains {
			// The user has typed in an exit code.
			return "", &Error{"The user has exited the program", 0}
		}

		return output, nil
	}

	return "", &Error{"Could not read the next line from standard input.", 1}
}

// Error specifies the whisper specific error type that can be tested against
type Error struct {
	text string // The message of the error
	code int    // The code of the error
}

func (err *Error) Error() string {
	return err.text
}

// Code returns the error code of the message to exit the program.
func (err *Error) Code() int {
	return err.code
}
