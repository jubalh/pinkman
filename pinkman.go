package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bobappleyard/readline"
	"github.com/codegangsta/cli"
	"github.com/wfreeman/pgn"
	"github.com/wfreeman/uci"
)

var b *pgn.Board

func run(cli *cli.Context) {
	var errmsg string
	var infomsg string
	var prompt string

	running := false

	uciPath := cli.GlobalString("path")

	engine, err := uci.NewEngine(uciPath)
	if err != nil {
		fmt.Println("Error: Could not start UCI engine from:", uciPath)
		fmt.Println(err)
		return
	}
	engine.SetOptions(uci.Options{
		Hash:    128,
		Ponder:  false,
		OwnBook: true,
		MultiPV: 4,
	})

	printWelcome()

	b = pgn.NewBoard()
	activePlayer := "white"

	for {
		if infomsg != "" {
			prompt += infomsg
		}
		if errmsg != "" {
			prompt += " \u25AB "
			prompt += errmsg
		}
		if running {
			prompt += " \u25AB " + activePlayer
		}
		prompt += "# "

		inputline, err := readline.String(prompt)
		prompt = ""
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error: ", err)
			break
		}
		switch inputline {
		case "start":
			infomsg = "Game started"
			running = true
		case "stop":
			infomsg = "Game stopped"
			running = false
		case "exit":
			engine.Close()
			return
		case "showfen":
			fmt.Println("FEN: ", b.String())
		default:
			if running {
				if len(inputline) >= 4 {
					err = b.MakeCoordMove(inputline)
					if err != nil && err != pgn.ErrUnknownMove {
						errmsg = err.Error()
						break
					}
				}
				if !cli.GlobalIsSet("no-engine") {
					engine.SetFEN(b.String())
					resultOps := uci.HighestDepthOnly
					results, err := engine.GoDepth(10, resultOps)
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println("Best move:", results.BestMove)
					err = b.MakeCoordMove(results.BestMove)
					break
				}
			}
		}
		fen := pgn.FENFromBoard(b)
		if fen.ToMove == pgn.White {
			activePlayer = "white"
		} else {
			activePlayer = "black"
		}
		if running {
			drawBoard(fen.FOR)
		}
		readline.AddHistory(inputline)
	}
}

func cmdShow(*cli.Context) {
}
func main() {
	app := cli.NewApp()

	app.Name = "pinkman"
	app.Usage = description
	app.Author = "Michael Vetter"
	app.Version = "0.0.1"
	app.Email = "g.bluehut@gmail.com"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "no-engine",
			Usage: "don't play against an engine. Can be used in case two people want to play on the same computer"},
		cli.StringFlag{
			Name:  "path",
			Value: "./stockfish",
			Usage: "set path to UCI engine. Standard is a executable named 'stockfish' in the same directory as the pinkman binary",
		},
	}

	app.Action = run

	app.Run(os.Args)
}
