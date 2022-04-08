package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/halvfigur/rhttp"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "rhttp",
		Usage: "HTTP tunnel client",
		Before: func(c *cli.Context) error {
			termChan := make(chan os.Signal, 1)
			signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

			c.App.Metadata["termChan"] = termChan
			return nil
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "port",
				Usage:    "local `PORT`",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "proxy",
				Usage:    "address and port of remote `PROXY`",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			client, err := rhttp.NewClient(c.Context, c.String("proxy"), c.Int("port"))
			if err != nil {
				return err
			}

			fmt.Println("Proxy ", client.ProxyAddr())
			<-termChan(c)
			return client.Close()
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func termChan(c *cli.Context) chan os.Signal {
	return c.App.Metadata["termChan"].(chan os.Signal)
}
