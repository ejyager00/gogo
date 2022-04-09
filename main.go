package main

import (
	"errors"
	"fmt"
)

type GoBoard = [][]int8
type GoGame struct {
	size              int
	board             GoBoard
	history           []string
	turn              int8
	consecutivePasses int
	piecesCaptured    [2]int
	komi              float32
}

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

func NewGame(sideLength int, komi float32) GoGame {
	game := GoGame{sideLength, newBoard(sideLength), nil, 1, 0, [2]int{0, 0}, komi}
	game.history = startHistory(game.board)
	return game
}

func newBoard(side int) GoBoard {
	board := make([][]int8, side)
	for i := range board {
		board[i] = make([]int8, side)
	}
	return board
}

func startHistory(board GoBoard) []string {
	return []string{fmt.Sprint(board)}
}

func (game GoGame) PrintBoard() {
	for i := range game.board {
		for j := range game.board[i] {
			if game.board[i][j] == 1 {
				fmt.Print("●")
			} else if game.board[i][j] == -1 {
				fmt.Print("○")
			} else {
				fmt.Print("+")
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func (game *GoGame) MakeMove(x, y int, pass, resign bool) (int8, error) {
	if pass {
		game.consecutivePasses++
		if game.consecutivePasses > 1 {
			scores := game.GetScores()
			if scores[0] > scores[1] {
				return 1, nil
			} else {
				return -1, nil
			}
		}
		game.turn *= -1
		return 0, nil
	}
	if resign {
		return game.turn * -1, nil
	}
	if game.board[y][x] != 0 {
		return 0, fmt.Errorf("there is already a piece at %d, %d", x, y)
	}
	boardCopy := make([][]int8, game.size)
	copy(boardCopy, game.board)
	boardCopy[y][x] = game.turn
	if y-1 >= 0 && boardCopy[y-1][x] == -1*game.turn {
		friends := getConnectedFriends(boardCopy, x, y-1, game.size, [][2]int{{x, y - 1}})
		if friendsSurrounded(boardCopy, friends, -1*game.turn, game.size) {
			for _, f := range friends {
				boardCopy[f[1]][f[0]] = 0
				if game.turn == 1 {
					game.piecesCaptured[0]++
				} else {
					game.piecesCaptured[1]++
				}
			}
		}
	}
	if y+1 < game.size && boardCopy[y+1][x] == -1*game.turn {
		friends := getConnectedFriends(boardCopy, x, y+1, game.size, [][2]int{{x, y + 1}})
		if friendsSurrounded(boardCopy, friends, -1*game.turn, game.size) {
			for _, f := range friends {
				boardCopy[f[1]][f[0]] = 0
				if game.turn == 1 {
					game.piecesCaptured[0]++
				} else {
					game.piecesCaptured[1]++
				}
			}
		}
	}
	if x-1 >= 0 && boardCopy[y][x-1] == -1*game.turn {
		friends := getConnectedFriends(boardCopy, x-1, y, game.size, [][2]int{{x - 1, y}})
		if friendsSurrounded(boardCopy, friends, -1*game.turn, game.size) {
			for _, f := range friends {
				boardCopy[f[1]][f[0]] = 0
				if game.turn == 1 {
					game.piecesCaptured[0]++
				} else {
					game.piecesCaptured[1]++
				}
			}
		}
	}
	if x+1 < game.size && boardCopy[y][x+1] == -1*game.turn {
		friends := getConnectedFriends(boardCopy, x+1, y, game.size, [][2]int{{x + 1, y}})
		if friendsSurrounded(boardCopy, friends, -1*game.turn, game.size) {
			for _, f := range friends {
				boardCopy[f[1]][f[0]] = 0
				if game.turn == 1 {
					game.piecesCaptured[0]++
				} else {
					game.piecesCaptured[1]++
				}
			}
		}
	}
	friends := getConnectedFriends(boardCopy, x, y, game.size, [][2]int{{x, y}})
	if friendsSurrounded(boardCopy, friends, game.turn, game.size) {
		for _, f := range friends {
			boardCopy[f[1]][f[0]] = 0
			if game.turn == 1 {
				game.piecesCaptured[1]++
			} else {
				game.piecesCaptured[0]++
			}
		}
	}
	if inHistory(boardCopy, &game.history) {
		return 0, errors.New("this move violates the ko rule")
	}
	copy(game.board, boardCopy)
	game.turn *= -1
	game.history = append(game.history, fmt.Sprint(game.board))
	game.consecutivePasses = 0
	return 0, nil
}

func (game GoGame) GetScores() [2]float32 {
	scores := [2]float32{0, 0}
	scores[1] += game.komi
	scores[0] += float32(game.piecesCaptured[0])
	scores[1] += float32(game.piecesCaptured[1])
	boardCopy := make([][]int8, game.size)
	copy(boardCopy, game.board)
	for i := 0; i < game.size; i++ {
		for j := 0; j < game.size; j++ {
			if boardCopy[i][j] == 0 {
				block := getConnectedFriends(boardCopy, j, i, game.size, [][2]int{{j, i}})
				if spaceEnclosed(boardCopy, block, 1, game.size) {
					for _, b := range block {
						boardCopy[b[1]][b[0]] = 1
					}
				} else if spaceEnclosed(boardCopy, block, -1, game.size) {
					for _, b := range block {
						boardCopy[b[1]][b[0]] = -1
					}
				}
			}
		}
	}
	for i := 0; i < game.size; i++ {
		for j := 0; j < game.size; j++ {
			if boardCopy[i][j] == 1 {
				scores[0]++
			} else if boardCopy[i][j] == -1 {
				scores[1]++
			}
		}
	}
	return scores
}

func getConnectedFriends(board GoBoard, x, y, side int, friends [][2]int) [][2]int {
	friendsCopy := make([][2]int, len(friends))
	copy(friendsCopy, friends)
	piece := board[y][x]
	if y-1 >= 0 && board[y-1][x] == piece {
		if !friendInSet(x, y-1, friendsCopy) {
			friendsCopy = append(friendsCopy, [2]int{x, y - 1})
			friendsCopy = getConnectedFriends(board, x, y-1, side, friendsCopy)
		}
	}
	if y+1 < side && board[y+1][x] == piece {
		if !friendInSet(x, y+1, friendsCopy) {
			friendsCopy = append(friendsCopy, [2]int{x, y + 1})
			friendsCopy = getConnectedFriends(board, x, y+1, side, friendsCopy)
		}
	}
	if x-1 >= 0 && board[y][x-1] == piece {
		if !friendInSet(x-1, y, friendsCopy) {
			friendsCopy = append(friendsCopy, [2]int{x - 1, y})
			friendsCopy = getConnectedFriends(board, x-1, y, side, friendsCopy)
		}
	}
	if x+1 < side && board[y][x+1] == piece {
		if !friendInSet(x+1, y, friendsCopy) {
			friendsCopy = append(friendsCopy, [2]int{x + 1, y})
			friendsCopy = getConnectedFriends(board, x+1, y, side, friendsCopy)
		}
	}
	return friendsCopy
}

func friendInSet(x, y int, friends [][2]int) bool {
	for i := range friends {
		if friends[i][0] == x && friends[i][1] == y {
			return true
		}
	}
	return false
}

func friendsSurrounded(board GoBoard, friends [][2]int, player int8, side int) bool {
	for _, f := range friends {
		if f[0]-1 >= 0 && board[f[1]][f[0]-1] == 0 {
			return false
		}
		if f[0]+1 < side && board[f[1]][f[0]+1] == 0 {
			return false
		}
		if f[1]-1 >= 0 && board[f[1]-1][f[0]] == 0 {
			return false
		}
		if f[1]+1 < side && board[f[1]+1][f[0]] == 0 {
			return false
		}
	}
	return true
}

func spaceEnclosed(board GoBoard, zeros [][2]int, player int8, side int) bool {
	for _, f := range zeros {
		if f[0]-1 >= 0 && board[f[1]][f[0]-1] == -1*player {
			return false
		}
		if f[0]+1 < side && board[f[1]][f[0]+1] == -1*player {
			return false
		}
		if f[1]-1 >= 0 && board[f[1]-1][f[0]] == -1*player {
			return false
		}
		if f[1]+1 < side && board[f[1]+1][f[0]] == -1*player {
			return false
		}
	}
	return true
}

func inHistory(board GoBoard, gameHistory *[]string) bool {
	b := fmt.Sprint(board)
	for _, x := range *gameHistory {
		if b == x {
			return true
		}
	}
	return false
}
