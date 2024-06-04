package main

import (
	"log"
	"os"

	kval "github.com/kval-access-language/kval-boltdb"
	"github.com/urfave/cli"
)

// appVersion is the current version of the app.
const appVersion = "2.0.3"

func main() {
	var file string
	var noValues bool
	var useMore bool

	app := cli.App{
		Name:    "bolter",
		Usage:   "view boltdb files interactively in your terminal",
		Version: appVersion,
		Authors: []cli.Author{
			{
				Name:  "Hasit Mistry",
				Email: "hasitnm@gmailcom",
			},
			{
				Name:  "vanillaiice",
				Email: "vanillaiice1@proton.me",
			},
		},
		Copyright: "(c) 2024 Hasit Mistry, vanillaiice",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file, f",
				Usage:       "load boltdb `FILE`",
				Destination: &file,
			},
			&cli.BoolFlag{
				Name:        "no-values",
				Usage:       "do not print values (use if values are huge and/or not printable)",
				Destination: &noValues,
			},
			&cli.BoolFlag{
				Name:        "more",
				Usage:       "use `more` to print all listings",
				Destination: &useMore,
			},
		},
		Action: func(c *cli.Context) (err error) {
			if file == "" {
				cli.ShowAppHelp(c)
				return
			}

			var formatter Formatter
			if useMore {
				formatter = &MoreWrapFormatter{
					formatter: formatter,
				}
			} else {
				formatter = &TableFormatter{
					noValues: noValues,
				}
			}

			if _, err = os.Stat(file); os.IsNotExist(err) {
				log.Fatal(err)
				return
			}

			i := Impl{fmt: formatter}
			i.initDB(file)

			defer kval.Disconnect(i.kb)

			i.readInput()

			return
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
