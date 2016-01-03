package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bobappleyard/readline"
	"github.com/codegangsta/cli"
	"github.com/freeeve/pgn"
	"github.com/jubalh/uci"
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

func parse_options(cli *cli.Context) {
	game.against_ai = !cli.GlobalIsSet("no-ai")
	game.engine_path = cli.GlobalString("path")

	if cli.GlobalIsSet("ai-white") {
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

func get_engine_move() string {
	game.engine.SetFEN(board.String())

	resultOps := uci.HighestDepthOnly
	results, err := game.engine.GoDepth(10, resultOps)
	if err != nil {
		fmt.Println("do something")
		fmt.Println(err)
	}
	return results.BestMove
}

func terminate(err error) {
	if err != nil {
		fmt.Println("erro: ", err)
		os.Exit(1)
	}
}

func make_turn(inputline string) string {
	// if AI move
	if game.active_player == game.ai_color {
		move := get_engine_move()
		board.MakeCoordMove(move)
		next_player()
	} else {
		game.engine.SetFEN(board.String())

		legal, err := game.engine.IsLegalMove(inputline)
		terminate(err)
		if legal {
			err = board.MakeCoordMove(inputline)
			if err != nil && err != pgn.ErrUnknownMove {
				return "Illegal Move"
			}
			next_player()
		} else {
			return "Illegal Move"
		}
	}
	return ""
}

func run(cli *cli.Context) {
	var errmsg string
	var infomsg string

	game.active = false
	parse_options(cli)

	err := launch_engine(cli)
	if err != nil {
		fmt.Println("Error: Could not start UCI engine from:", game.engine_path)
		return
	}

	printWelcome()

	for {
		inputline, err := get_readline_prompt(infomsg, errmsg)
		if err == io.EOF {
			return
		}
		terminate(err)
		errmsg = ""

		switch inputline {
		// commands
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
			return
		case "showfen":
			fmt.Println("FEN: ", board.String())
			break
		// moves
		default:
			if game.active {
				if game.ai_color == "white" && game.active_player == game.ai_color {
					make_turn("")
				}
				if len(inputline) >= 4 {
					errmsg = make_turn(inputline)
					if game.against_ai && len(errmsg) < 1 {
						make_turn("")
					}
				}
			}

			if game.active {
				draw_board()
			}

			readline.AddHistory(inputline)
		}
	}
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

	app.Run(os.Args)
}
