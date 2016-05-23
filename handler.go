package whisper

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Strings to specify exit queries on input.
const (
	EXIT = "exit"
	QUIT = "quit"
)

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
