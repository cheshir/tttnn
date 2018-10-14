package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var xml2csvCommand = cli.Command{
	Name:  "xml2csv",
	Usage: "Parses XML with gomoku games and saves them to CSV file",
	Action: func(c *cli.Context) error {
		from := c.String("from")
		to := c.String("to")

		if from == "" || to == "" {
			return fmt.Errorf("you must specify `from` and `to` files")
		}

		// Open files.
		fromFile, err := os.OpenFile(from, os.O_RDONLY, 0666)
		if err != nil {
			return errors.Wrapf(err, "failed to open file %s", fromFile)
		}
		defer fromFile.Close()

		toFile, err := os.Create(to)
		if err != nil {
			return errors.Wrapf(err, "failed to open file %s", toFile)
		}
		defer toFile.Close()

		// Decode.
		destination := csv.NewWriter(toFile)

		var games renjunetGames
		if err := xml.NewDecoder(NewValidUTF8Reader(fromFile)).Decode(&games); err != nil {
			return errors.Wrap(err, "failed to decode xml file")
		}

		log.Printf("Decoded %d games\n", len(games.Games))

		for _, game := range games.Games {
			if len(game.Moves) > 0 {
				destination.Write([]string{game.BlackResult, game.Moves[0]})
			}
		}

		destination.Flush()
		if err := destination.Error(); err != nil {
			return errors.Wrap(err, "failed to write csv file")
		}

		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "from",
			Usage: "XML `FILE` that we're going to parse",
		},
		cli.StringFlag{
			Name:  "to",
			Usage: "CSV `FILE` where we're going to save result",
		},
	},
}

type renjunetGames struct {
	XMLName xml.Name       `xml:"database"`
	Games   []renjunetGame `xml:"games>game"`
}

type renjunetGame struct {
	BlackResult string   `xml:"bresult,attr"`
	Moves       []string `xml:"move"`
}
