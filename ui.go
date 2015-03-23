package main

import "fmt"

func drawColumnsIndex() {
	fmt.Print("  ")
	for i := 'A'; i <= 'H'; i++ {
		fmt.Print("  ", string(i), " ")
	}
	fmt.Println()
}

func drawRowSeperator() {
	fmt.Print("   ")
	for i := 0; i < 8; i++ {
		fmt.Print("____")
	}
	fmt.Println()
}

func drawBoard(sFOR string) {
	symbol := ' '
	row := 8

	drawColumnsIndex()
	drawRowSeperator()
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
	drawColumnsIndex()
	drawRowSeperator()
}
