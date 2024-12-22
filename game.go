package main

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"time"
)

// MakeMove attempts to set a value at a specific cell.
func (g *Game) MakeMove(row, col, value int) error {

	Log.Debug().Str("game_id", g.ID).Int("row", row).Int("col", col).Int("value", value).Msg("validating move")

	var err error

	// Ensure value is valid
	if value < 1 || value > 9 {
		err = errors.New("value must be between 1 and 9")
	}

	// Ensure cell is not fixed
	if g.Board[row][col].Fixed {
		err = errors.New("desired cell is fixed")
	}

	// Ensure move doesn't violate Sudoku rules
	if !isValidMove(g.Board, row, col, value) {
		err = errors.New("move violates Sudoku rules")
	}

	if err != nil {
		Log.Error().Str("game_id", g.ID).Int("row", row).Int("col", col).Int("value", value).Err(err).Msg("failed to make move")
		return err
	}

	// Record the previous value
	previousValue := g.Board[row][col].Value

	// Make the move
	g.Board[row][col].Value = value

	move := Move{
		Row:           row,
		Col:           col,
		PreviousValue: previousValue,
		CurrentValue:  value,
		Timestamp:     time.Now().Unix(),
	}

	// Add the move to the history
	g.MoveHistory = append(g.MoveHistory, move)

	Log.Info().Str("game_id", g.ID).Int("row", row).Int("col", col).Int("value", value).Msg("move successful")

	return nil
}

// RemoveMove clears the value of a cell if it's not fixed.
func (g *Game) RemoveMove(row, col int) error {

	if g.Board[row][col].Fixed {
		err := errors.New("cannot remove value from a fixed cell")
		Log.Error().Str("game_id", g.ID).Int("row", row).Int("col", col).Int("value", 0).Err(err).Msg("failed to remove move")
		return err
	}

	previousValue := g.Board[row][col].Value

	// Clear the cell
	g.Board[row][col].Value = 0

	// Add the removal action to history
	move := Move{
		Row:           row,
		Col:           col,
		PreviousValue: previousValue,
		CurrentValue:  0,
		Timestamp:     time.Now().Unix(),
	}

	g.MoveHistory = append(g.MoveHistory, move)

	Log.Info().Str("game_id", g.ID).Int("row", row).Int("col", col).Int("value", 0).Msg("move removed with successful")

	return nil
}

// GetHint provides a hint for the current game by finding an empty cell and suggesting a valid value.
func (g Game) GetHint() (*Hint, error) {
	// Find the first empty cell and suggest a valid value
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if g.Board[i][j].Value == 0 {
				for val := 1; val <= 9; val++ {
					if isValidMove(g.Board, i, j, val) {
						return &Hint{Row: i, Col: j, Value: val}, nil
					}
				}
			}
		}
	}

	return nil, errors.New("no hints available")
}

// SolveGame solves the Sudoku puzzle using a backtracking algorithm.
func (g *Game) SolveGame() error {

	var tries *uint64
	tries = new(uint64)
	*tries = 0

	if !solveBoard(&g.Board, tries) {
		return errors.New("unable to solve the puzzle")
	}

	log.Info().Str("game_id", g.ID).Uint64("tries", *tries).Msg("game solved")

	return nil
}

func (b Board) String() string {
	var result string

	for i, row := range b {
		for j, cell := range row {
			if cell.Value == 0 {
				result += ". "
			} else {
				result += fmt.Sprintf("%d ", cell.Value)
			}
			if j%3 == 2 && j < 8 {
				result += "| "
			}
		}
		result += "\n"
		if i%3 == 2 && i < 8 {
			result += "------+-------+------\n"
		}
	}
	return result
}

func (g Game) String() string {
	var result string
	result += "Game ID: " + g.ID + "\n"
	result += "Difficulty: " + g.Difficulty.String() + "\n"
	result += "Total Moves: " + fmt.Sprintf("%d", len(g.MoveHistory)) + "\n"
	result += "Board:\n"
	result += g.Board.String()
	return result
}
