package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/bbengfort/whisper"
	"github.com/codegangsta/cli"
	"github.com/joho/godotenv"
)

func main() {

	// Load the .env file if it exists
	godotenv.Load()

	// Instantiate the command line application.
	app := cli.NewApp()
	app.Name = "whisper"
	app.Usage = "text based P2P whispernet client"
	app.Version = whisper.Version
	app.Author = "Benjamin Bengfort"
	app.Email = "bengfort@cs.umd.edu"
	app.EnableBashCompletion = true
	app.Action = beginWhispering

	// Set the flags and options on the app
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "name",
			Value:  getDefaultName(),
			Usage:  "identify yourself on the whispernet",
			EnvVar: "WHISPER_NAME",
		},
	}

	// Run the command line application
	app.Run(os.Args)
}

// Begins listening for command line input and writes messages.
func beginWhispering(ctx *cli.Context) error {

	// Run the interrupt handler
	handleInterrupt()

	// Collect the name from the command line input
	name := ctx.String("name")

	reader := whisper.NewInputHandler(">")
	for {
		body, err := reader.ReadLine()
		if err != nil {
			return cli.NewExitError(err.Error(), err.Code())
		}

		// Construct the message to write as output.
		msg := whisper.NewMessage(body, name)
		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(msg); err != nil {
			return cli.NewExitError(err.Error(), 3)
		}
	}
}

// Gets the default name from the environment
func getDefaultName() string {
	// Try for the hostname first
	name, err := os.Hostname()
	if err == nil {
		return name
	}

	return "unknown"
}

// Create the interrupt handler
func handleInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println("Shutting down!")
			os.Exit(0)
		}
	}()
}
