// Pinkman is a command line chess interface to UCI chess engines
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bobappleyard/readline"
	"github.com/freeeve/pgn"
	"github.com/jubalh/uci"
	"gopkg.in/urfave/cli.v1"
)

type Game struct {
	against_ai    bool
	engine_path   string
	engine        *uci.Engine
	active_player string
	active        bool
	ai_color      string
}

var board *pgn.Board
var game Game

func next_player() {
	if game.active_player == "white" {
		game.active_player = "black"
	} else {
		game.active_player = "white"
	}
}

func parse_options(ctx *cli.Context) {
	game.against_ai = !ctx.GlobalIsSet("no-ai")
	game.engine_path = ctx.GlobalString("path")

	if ctx.GlobalIsSet("ai-white") {
		game.ai_color = "white"
	} else {
		game.ai_color = "black"
	}
}

func launch_engine(cli *cli.Context) error {
	engine, err := uci.NewEngine(game.engine_path)
	if err != nil {
		return err
	}
	engine.SetOptions(uci.Options{
		Hash:    128,
		Ponder:  false,
		OwnBook: true,
		MultiPV: 4,
	})
	game.engine = engine
	return nil
}

func get_readline_prompt(infomsg, errmsg string) (string, error) {
	var prompt string

	if infomsg != "" {
		prompt += infomsg
	}
	if errmsg != "" {
		prompt += " \u25AB "
		prompt += errmsg
	}
	if game.active {
		prompt += " \u25AB " + game.active_player
	}
	prompt += "# "

	return readline.String(prompt)
}

func draw_board() {
	fen := pgn.FENFromBoard(board)
	drawBoard(fen.FOR)
}

func get_engine_move() (string, error) {
	game.engine.SetFEN(board.String())

	resultOps := uci.HighestDepthOnly
	results, err := game.engine.GoDepth(10, resultOps)
	if err != nil {
		return "", err
	}

	return results.BestMove, nil
}

func make_turn(inputline string) (string, error) {
	// if AI move
	if game.active_player == game.ai_color {
		move, err := get_engine_move()
		if err != nil {
			return "", err
		}
		board.MakeCoordMove(move)
		next_player()
	} else {
		game.engine.SetFEN(board.String())

		legal, err := game.engine.IsLegalMove(inputline)
		if err != nil {
			return "", err
		}
		if legal {
			err = board.MakeCoordMove(inputline)
			if err != nil && err != pgn.ErrUnknownMove {
				return "Illegal move", nil
			}
			next_player()
		} else {
			return "Illegal move", nil
		}
	}
	return "", nil
}

func run(ctx *cli.Context) error {
	var errmsg string
	var infomsg string

	game.active = false
	parse_options(ctx)

	err := launch_engine(ctx)
	if err != nil {
		return cli.NewExitError("Error: Could not start UCI engine from: "+game.engine_path, 1)
	}

	printWelcome()

	for {
		inputline, err := get_readline_prompt(infomsg, errmsg)
		if err == io.EOF {
			return cli.NewExitError("Error: "+err.Error(), 1)
		}
		errmsg = ""

		switch inputline {
		// commands
		case "help":
			printInGameHelp()
			break
		case "start":
			board = pgn.NewBoard()
			infomsg = "Game started"
			game.active = true
			game.active_player = "white"
			draw_board()
			break
		case "stop":
			infomsg = "Game stopped"
			game.active = false
			break
		case "exit":
			if game.engine != nil {
				game.engine.Close()
			}
			return nil
		case "showfen":
			fmt.Println("FEN: ", board.String())
			break
		// moves
		default:
			if game.active {
				if game.ai_color == "white" && game.active_player == game.ai_color {
					_, err := make_turn("")
					if err != nil {
						return cli.NewExitError("Error: "+err.Error(), 1)
					}
				}
				if len(inputline) >= 4 {
					errmsg, err = make_turn(inputline)
					if err != nil {
						return cli.NewExitError("Error: "+err.Error(), 1)
					}
					if game.against_ai && len(errmsg) < 1 {
						_, err := make_turn("")
						if err != nil {
							return cli.NewExitError("Error: "+err.Error(), 1)
						}
					}
				}
			}

			if game.active {
				draw_board()
			}

			readline.AddHistory(inputline)
		}
	}
	return nil
}

func main() {
	app := cli.NewApp()

	app.Name = "pinkman"
	app.Usage = description
	app.Author = "Michael Vetter"
	app.Version = "0.1"
	app.Email = "jubalh@openmailbox.org"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "no-ai",
			Usage: "don't play against the engine. Can be used in case two people want to play on the same computer"},
		cli.BoolFlag{
			Name:  "ai-white",
			Usage: "let AI play white"},
		cli.StringFlag{
			Name:  "path",
			Value: "./stockfish",
			Usage: "set path to UCI engine. Standard is a executable named 'stockfish' in the same directory as the pinkman binary",
		},
	}

	app.Action = run

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
}
