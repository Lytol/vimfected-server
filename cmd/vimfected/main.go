package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Lytol/vimfected-server/game"
	"github.com/Lytol/vimfected-server/server"
)

func main() {
	app := &cli.App{
		Name:  "server",
		Usage: "run a vimfected server instance",
		Action: func(*cli.Context) error {
			g, err := game.New()
			if err != nil {
				return err
			}

			s, err := server.NewServer(g)
			if err != nil {
				return err
			}
			return s.Run()
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
