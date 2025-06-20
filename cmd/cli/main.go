package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "greet1",
		Usage: "fight the loneliness!",
		Action: func(c *cli.Context) error {
			fmt.Println("you just typed greet command")
			fmt.Println(c.Args().First())
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "first-command",
				Usage: "this is the first command",
				Action: func(c *cli.Context) error {
					fmt.Println("you just typed the first command")
					fmt.Println(c.Args().First())
					return nil
				},
			},
			{
				Name:  "second-command",
				Usage: "this is the second command",
				Action: func(c *cli.Context) error {
					fmt.Println("you just typed the second command")
					fmt.Println(c.Args().First())
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
