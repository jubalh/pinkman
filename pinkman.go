package main

import (
	"fmt"
	"github.com/bobappleyard/readline"
	"github.com/wfreeman/pgn"
	"io"
)

var b *pgn.Board

func main() {
	var errmsg string
	var infomsg string
	var prompt string
	running := false
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

func _draw_columns_index() {
	fmt.Print("  ")
	for i := 'A'; i <= 'H'; i++ {
		fmt.Print("  ", string(i), " ")
	}
	fmt.Println()
}

func _draw_row_seperator() {
	fmt.Print("   ")
	for i := 0; i < 8; i++ {
		fmt.Print("____")
	}
	fmt.Println()
}

func drawBoard(sFOR string) {
	symbol := ' '
	row := 8

	_draw_columns_index()
	_draw_row_seperator()
	fmt.Print(row, "| ")
	row--
	for _, c := range sFOR {
		switch {
		case c == 'p':
			symbol = '♟'
		case c == 'r':
			symbol = '♜'
		case c == 'n':
			symbol = '♞'
		case c == 'b':
			symbol = '♝'
		case c == 'q':
			symbol = '♛'
		case c == 'k':
			symbol = '♚'
		case c == 'P':
			symbol = '♙'
		case c == 'R':
			symbol = '♖'
		case c == 'N':
			symbol = '♘'
		case c == 'B':
			symbol = '♗'
		case c == 'Q':
			symbol = '♕'
		case c == 'K':
			symbol = '♔'
		case c == '/':
			fmt.Println()
			fmt.Print(row, "| ")
			row--
		case c >= '1' && c <= '8':
			for i := 0; i < int(c-'0'); i++ {
				fmt.Print("[  ]")
			}
		}
		if symbol != ' ' {
			fmt.Print("[", string(symbol), " ]")
			symbol = ' '
		}
	}
	fmt.Println()
	_draw_columns_index()
	_draw_row_seperator()
}
