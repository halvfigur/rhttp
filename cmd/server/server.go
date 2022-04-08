package main

import (
	"fmt"
	"os"

	"github.com/halvfigur/rhttp"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "rhttp",
		Usage: "HTTP tunnel server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "port",
				Usage:    "local `PORT`",
				Value:    80,
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			return rhttp.NewServer(c.Int("port"))
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
