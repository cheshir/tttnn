package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/cheshir/tttnn/game"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var filterFullGamesCommand = cli.Command{
	Name:  "filter-full-games",
	Usage: "Filters prepared raw games in CSV format",
	Action: func(c *cli.Context) error {
		from := c.String("from")
		to := c.String("to")

		if from == "" || to == "" {
			return fmt.Errorf("you must specify `from` and `to` files")
		}

		// Open files.
		fromFile, err := os.OpenFile(from, os.O_RDONLY, 0666)
		if err != nil {
			return errors.Wrapf(err, "failed to open file %s", from)
		}
		defer fromFile.Close()

		toFile, err := os.Create(to)
		if err != nil {
			return errors.Wrapf(err, "failed to open file %s", to)
		}
		defer toFile.Close()

		// Processing.
		source := csv.NewReader(fromFile)
		destination := csv.NewWriter(toFile)

		var finishedCount, abortedCount, wrongResultCount, i int
		for {
			i++
			record, err := source.Read()
			if err != nil {
				if err == io.EOF {
					break
				}

				return errors.Wrap(err, "failed to read from csv file")
			}

			actualGameResult, err := simulateGame(record[1])
			if err != nil {
				wrongResultCount++
				log.Println(errors.Wrapf(err, "game simulation failed. game: %s", record[1]))

				continue
			}

			if actualGameResult == game.Undefined {
				abortedCount++

				continue
			}

			expectedResult := convertExternalResultToInternal(record[0])
			if actualGameResult != expectedResult {
				wrongResultCount++

				fixed := convertInternalResultToExternal(actualGameResult)
				if fixed != "" {
					record[0] = fixed
					fmt.Print("[fixed] ")
				}

				log.Printf("simulated result %v is not equal to expected %v. game: %s\n", actualGameResult, expectedResult, record[1])
			}

			destination.Write(record)
			finishedCount++

			if i%10000 == 0 {
				printStats(i, finishedCount, abortedCount, wrongResultCount)
			}
		}

		printStats(i, finishedCount, abortedCount, wrongResultCount)

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

func simulateGame(moves string) (result game.Result, err error) {
	table := game.New()

	for _, move := range strings.Fields(moves) {
		result, err = table.Move(move)
		if err != nil {
			return
		}
	}

	return
}

func convertExternalResultToInternal(result string) game.Result {
	switch result {
	case "1":
		return game.BlackWin
	case "0.5":
		return game.Draw
	case "0":
		return game.BlackLose
	}

	return game.Invalid
}

func convertInternalResultToExternal(result game.Result) string {
	switch result {
	case game.BlackWin:
		return "1"
	case game.BlackLose:
		return "0"
	case game.Draw:
		return "0.5"
	}

	return ""
}

func printStats(total, finished, aborted, failed int) {
	fmt.Printf(`
total:    %d
finished: %d
aborted:  %d
failed:   %d
`,
		total, finished, aborted, failed)
}
