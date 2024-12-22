package main

import (
	"fmt"
	"time"
)

func solveBoard(board *Board, tries *uint64) bool {

	Log.Debug().Uint64("tries", *tries).Msg("Trying to solve board")

	*tries = *tries + 1

	row, col, empty := findEmptyCell(*board)

	if !empty {
		return true // No empty cells, puzzle solved
	}

	for val := 1; val <= 9; val++ {

		if isValidMove(*board, row, col, val) {
			board[row][col].Value = val

			// Recursively attempt to solve
			if solveBoard(board, tries) {
				return true
			}

			// Backtrack
			board[row][col].Value = 0
		}
	}

	return false
}

func findEmptyCell(board Board) (int, int, bool) {

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j].Value == 0 {
				return i, j, true
			}
		}
	}

	return 0, 0, false
}

func isValidMove(board Board, row, col, value int) bool {
	// Check row
	for i := 0; i < 9; i++ {
		if board[row][i].Value == value {
			return false
		}
	}

	// Check column
	for i := 0; i < 9; i++ {
		if board[i][col].Value == value {
			return false
		}
	}

	// Check 3x3 subgrid
	startRow, startCol := (row/3)*3, (col/3)*3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[startRow+i][startCol+j].Value == value {
				return false
			}
		}
	}

	return true
}

func generateBoard(opts BoardGeneratorOpts) (Board, error) {

	start := time.Now()
	Log.Info().Str("difficulty", opts.Difficulty.String()).Str("board_id", opts.BoardId).Msg("Generating board")

	var cellsToRemove int

	switch opts.Difficulty {
	case Easy:
		cellsToRemove = 38
	case Medium:
		cellsToRemove = 46
	case Hard:
		cellsToRemove = 49
	case Expert:
		cellsToRemove = 52
	case Master:
		cellsToRemove = 55
	case Sadistic:
		cellsToRemove = 58
	default:
		Log.Error().Int("difficulty", int(opts.Difficulty)).Str("board_id", opts.BoardId).Msg("Invalid difficulty level")
		return Board{}, fmt.Errorf("invalid difficulty level: %d", opts.Difficulty)
	}

	var board Board
	var tries *uint64 = new(uint64)
	var solved bool = false

	*tries = 0

	for !solved {

		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				board[i][j].Value = 0
				board[i][j].Fixed = false
			}
		}

		Log.Debug().Uint64("tries", *tries).Str("board_id", opts.BoardId).Msg("Attempting to generate a valid board")

		solved = solveBoard(&board, tries)

		if !solved {
			Log.Debug().Uint64("tries", *tries).Str("board_id", opts.BoardId).Msg("Failed to generate a valid board")
		}
	}

	Log.Debug().Uint64("tries", *tries).Str("board_id", opts.BoardId).Msg("Generated a valid board, removing cells")

	source := opts.SeedSourceGenerate()

	Log.Debug().Str("board_id", opts.BoardId).Msg("Generated seed source")

	for i := 0; i < cellsToRemove; i++ {
		row, col := source.Intn(9), source.Intn(9)

		for board[row][col].Value == 0 {
			row, col = source.Intn(9), source.Intn(9)
		}

		board[row][col].Value = 0
	}

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j].Value != 0 {
				board[i][j].Fixed = true
			}
		}
	}

	duration := time.Since(start)

	Log.Info().Str("duration", duration.String()).Str("board_id", opts.BoardId).Msg("Generated board")

	return board, nil
}

func copyBoard(board Board) Board {
	var copiedBoard Board

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			copiedBoard[i][j].Value = board[i][j].Value
			copiedBoard[i][j].Fixed = board[i][j].Fixed
		}
	}

	return copiedBoard
}
