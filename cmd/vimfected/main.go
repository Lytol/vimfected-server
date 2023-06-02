package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "server",
		Usage: "run a vimfected server instance",
		Action: func(*cli.Context) error {
			fmt.Printf("Running vimfected...\n")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
