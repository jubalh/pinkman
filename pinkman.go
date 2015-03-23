package main

import (
	"fmt"
	"io"

	"github.com/bobappleyard/readline"
	"github.com/wfreeman/pgn"
	"github.com/wfreeman/uci"
)

var b *pgn.Board
var stockfishPath = "/home/sb/stockfish-6-linux/stockfish-6-linux/Linux/stockfish_6_x64"

func main() {
	var errmsg string
	var infomsg string
	var prompt string

	running := false
	engine, err := uci.NewEngine(stockfishPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	engine.SetOptions(uci.Options{
		Hash:    128,
		Ponder:  false,
		OwnBook: true,
		MultiPV: 4,
	})

	fmt.Println("*** pinkman ***")
	fmt.Println("the totally kafkaesque chess game")
	fmt.Println()

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
			return
		case "showfen":
			fmt.Println("FEN: ", b.String())
		default:
			if running {
				if activePlayer == "black" {
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
				if len(inputline) >= 4 {
					err = b.MakeCoordMove(inputline)
					if err != nil && err != pgn.ErrUnknownMove {
						errmsg = err.Error()
						break
					}
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
