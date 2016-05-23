package main

import (
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
		cli.StringFlag{
			Name:  "address, A",
			Value: "localhost:3264",
			Usage: "specify a whisper server to connect to",
		},
	}

	// Run the command line application
	app.Run(os.Args)
}

// Begins listening for command line input and writes messages.
func beginWhispering(ctx *cli.Context) error {

	// Run the interrupt handler
	handleInterrupt()

	// Create the Client
	client := whisper.NewClient(ctx.String("name"))

	// Connect to the Server and run the application
	err := client.Connect(ctx.String("address"))
	if err != nil {
		return cli.NewExitError(err.Error(), err.Code())
	}

	err = client.Run()
	client.Close()

	// Return the exit error value
	return cli.NewExitError(err.Error(), err.Code())
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
