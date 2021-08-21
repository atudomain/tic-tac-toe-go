package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type iBoard interface {
	move(player int, coordinates [2]int)
	getRowSums() [3]int
	getColumnSums() [3]int
	getDiagonalSums() [2]int
	getFields() [3][3]int
	asses(player int) int
	copy() iBoard
	print()
}

type Board struct {
	fields       [3][3]int
	rowSums      [3]int
	columnSums   [3]int
	diagonalSums [2]int
	turn         int
}

func newBoard() iBoard {
	return &Board{
		fields:       [3][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		rowSums:      [3]int{0, 0, 0},
		columnSums:   [3]int{0, 0, 0},
		diagonalSums: [2]int{0, 0},
		turn:         0,
	}
}

func (b *Board) copy() iBoard {
	fields := b.fields
	rowSums := b.rowSums
	columnSums := b.columnSums
	diagonalSums := b.diagonalSums
	turn := b.turn
	return &Board{
		fields:       fields,
		rowSums:      rowSums,
		columnSums:   columnSums,
		diagonalSums: diagonalSums,
		turn:         turn,
	}
}

func (b *Board) getFields() [3][3]int {
	return b.fields
}

func (b *Board) move(player int, coordinates [2]int) {
	b.fields[coordinates[0]][coordinates[1]] = player
	b.turn++
	b.rowSums[coordinates[0]] = b.rowSums[coordinates[0]] + player
	b.columnSums[coordinates[1]] = b.columnSums[coordinates[1]] + player
	if coordinates[0] == coordinates[1] {
		b.diagonalSums[0] = b.diagonalSums[0] + player
	}
	for i := 0; i < 3; i++ {
		if coordinates[0] == (2-i) && coordinates[1] == (0+i) {
			b.diagonalSums[1] = b.diagonalSums[1] + player
		}
	}
}

func (b *Board) getRowSums() [3]int {
	return b.rowSums
}

func (b *Board) getColumnSums() [3]int {
	return b.columnSums
}

func (b *Board) getDiagonalSums() [2]int {
	return b.diagonalSums
}

func (b *Board) asses(player int) int {
	// Check for direct win position
	for i := 0; i < 3; i++ {
		if b.rowSums[i] == player*3 {
			return 100
		}
		if b.columnSums[i] == player*3 {
			return 100
		}
	}
	for i := 0; i < 2; i++ {
		if b.diagonalSums[i] == player*3 {
			return 100
		}
	}
	return 0
}

func (b *Board) print() {
	var character string
	fmt.Println()
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			switch b.fields[i][j] {
			case -1:
				character = "x"
			case 1:
				character = "o"
			default:
				character = "_"
			}
			fmt.Printf("%s ", character)
		}
		fmt.Println()
	}
	fmt.Println()
}

type Move struct {
	score       int
	coordinates [2]int
}

type ParsedMove struct {
	player      int
	coordinates [2]int
}

type UpdatedMove struct {
	updatedScore int
	move         Move
}

func main() {
	// Example usage:
	//
	// ./engine 1 0 0 0 0 -1 0 0 0 1
	//
	// First input integer is player that searches for move (-1 for x or 1 for o).
	// The rest of integers are fields on the board where
	// -1 stands for x and 1 for o (0 for empty field).
	// They are read left to right row by row from top to bottom.
	//
	// Example output:
	//
	// 0 0 0
	//
	// First integer is score for the best move.
	// If it is positive, moving player is going to win.
	// If negative, moving player is going to lose.
	// The bigger absolute value up to 100, the outcome is closer regarding depth.
	// If 0, the game is going to be draw.
	//
	// The next to integers are coordinates for the best move.
	//
	board := newBoard()

	player, _ := strconv.Atoi(os.Args[1])
	parsedMoves := parseMoves(os.Args[2:])

	for i := 0; i < len(parsedMoves); i++ {
		board.move(parsedMoves[i].player, parsedMoves[i].coordinates)
	}

	maxScore := -999
	availableCoordinates := findAvailableCoordinates(board)

	availableMoves := make([]Move, 0, 9)
	for i := 0; i < len(availableCoordinates); i++ {
		score := calculateScore(player, player, availableCoordinates[i], board, 0)
		availableMoves = append(availableMoves, Move{score, availableCoordinates[i]})
		if score > maxScore {
			maxScore = score
		}
	}

	potentialMoves := make([]Move, 0, 9)
	for i := 0; i < len(availableMoves); i++ {
		if availableMoves[i].score == maxScore {
			potentialMoves = append(potentialMoves, availableMoves[i])
		}
	}

	bestMoves := calculateBestMoves(player, board, potentialMoves)

	finalMove := calculateRandomMove(bestMoves)

	fmt.Printf("%d %d %d\n", finalMove.score, finalMove.coordinates[0], finalMove.coordinates[1])
}

func parseMoves(arguments []string) []ParsedMove {
	parsedMoves := make([]ParsedMove, 0, 9)
	for i := 0; i < len(arguments); i++ {
		player, _ := strconv.Atoi(arguments[i])
		if player != 0 {
			x := i / 3
			y := i % 3
			parsedMoves = append(parsedMoves, ParsedMove{player, [2]int{x, y}})
		}
	}
	return parsedMoves
}

func calculateBestMoves(player int, board iBoard, moves []Move) []Move {
	updatedMoves := make([]UpdatedMove, 0, 9)
	maxScore := -999
	for i := 0; i < len(moves); i++ {
		currentScore := moves[i].score
		boardCopy := board.copy()
		boardCopy.move(player, moves[i].coordinates)
		rowSums := boardCopy.getRowSums()
		columnSums := boardCopy.getColumnSums()
		diagonalSums := boardCopy.getDiagonalSums()
		for i := 0; i < 3; i++ {
			if rowSums[i] > 0 {
				currentScore = currentScore + rowSums[i]
			}
			if columnSums[i] > 0 {
				currentScore = currentScore + columnSums[i]
			}
		}
		for i := 0; i < 2; i++ {
			if diagonalSums[i] > 0 {
				currentScore = currentScore + diagonalSums[i]
			}
		}
		if currentScore > maxScore {
			maxScore = currentScore
		}
		updatedMoves = append(updatedMoves, UpdatedMove{currentScore, moves[i]})
	}
	resultMoves := make([]Move, 0, 9)
	for i := 0; i < len(updatedMoves); i++ {
		if updatedMoves[i].updatedScore == maxScore {
			resultMoves = append(resultMoves, updatedMoves[i].move)
		}
	}
	return resultMoves
}

func calculateRandomMove(moves []Move) Move {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(moves))
	return moves[randomIndex]
}

func calculateScore(player int, currentPlayer int, coordinates [2]int, board iBoard, depth int) int {
	boardCopy := board.copy()
	// Should get score after coordinates move
	boardCopy.move(currentPlayer, coordinates)
	score := boardCopy.asses(currentPlayer)
	if score != 0 {
		if player == currentPlayer {
			return score - depth
		} else {
			return (-1)*score + depth
		}
	}
	availableCoordinates := findAvailableCoordinates(boardCopy)
	if len(availableCoordinates) == 0 {
		return 0
	}
	// If game not finished, should check all possible moves
	scores := make([]int, 0, 9)
	for i := 0; i < len(availableCoordinates); i++ {
		score := calculateScore(player, (-1)*currentPlayer, availableCoordinates[i], boardCopy, depth+1)
		scores = append(scores, score)
	}
	// Depending on player, maximize or minimize score
	if player == currentPlayer {
		return min(scores)
	} else {
		return max(scores)
	}
}

func min(slice []int) int {
	smallest := slice[0]
	for i := 0; i < len(slice); i++ {
		if slice[i] < smallest {
			smallest = slice[i]
		}
	}
	return smallest
}

func max(slice []int) int {
	largest := slice[0]
	for i := 0; i < len(slice); i++ {
		if slice[i] > largest {
			largest = slice[i]
		}
	}
	return largest
}

func findAvailableCoordinates(board iBoard) [][2]int {
	coordinates := make([][2]int, 0, 9)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board.getFields()[i][j] == 0 {
				coordinates = append(coordinates, [2]int{i, j})
			}
		}
	}
	return coordinates
}
