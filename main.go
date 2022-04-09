package main

import (
	"fmt"
)

func main() {
	fmt.Println("The board is x, y indexed from the top left corner starting with 0.")
	fmt.Println("Input moves with two integers separated by a space representing x and y respectively.")
	fmt.Println("To pass, input -1 and to resign input -2.")

	var s int
	var k float32
	fmt.Print("Input board size and komi (white advantage) separated by a space: ")
	fmt.Scanf("%d %f", &s, &k)

	game := NewGame(s, k)

	gameover := false
	for !gameover {
		game.PrintBoard()
		if game.turn == 1 {
			fmt.Print("Enter black's move: ")
		} else {
			fmt.Print("Enter white's move: ")
		}
		var move_x, move_y int
		fmt.Scanf("%d %d\n", &move_x, &move_y)
		var err error
		var victory int8
		if move_x == -1 {
			victory, err = game.MakeMove(0, 0, true, false)
		} else if move_x == -2 {
			victory, err = game.MakeMove(0, 0, false, true)
		} else {
			victory, err = game.MakeMove(move_x, move_y, false, false)
		}
		if err != nil {
			fmt.Println(err)
		}
		if victory != 0 {
			gameover = true
			scores := game.GetScores()
			if victory > 0 {
				fmt.Printf("Black wins, %.1f to %.1f\n", scores[0], scores[1])
			} else {
				fmt.Printf("White wins, %.1f to %.1f\n", scores[1], scores[0])
			}
		}
	}
}
