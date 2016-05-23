package whisper

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
