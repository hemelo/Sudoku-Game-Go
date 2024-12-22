package main

import (
	"Sudoku-Solver/pkg/logo"
	"fmt"
	"github.com/inancgumus/screen"
	"github.com/manifoldco/promptui"
	"strconv"
	"time"
)

const (
	CreateGameAction      = "Create Game"
	SelectGameAction      = "Select Game"
	ExitAction            = "Exit"
	ShowBoardAction       = "Show Board"
	MakeMoveAction        = "Make Move"
	RemoveMoveAction      = "Remove Move"
	GetHintAction         = "Get Hint"
	SolveAction           = "Solve"
	BackAction            = "Back to Main Menu"
	ShowGameHistoryAction = "Show Game MoveHistory"
)

func NewClientCLI(opts ClientCLIOpts) *ClientCLI {
	return &ClientCLI{
		defaultDifficulty: opts.defaultDifficulty,
		gameController:    opts.gameController,
		gameManager:       opts.gameManager,
		currentGame:       nil,
		currentGameID:     "",
	}
}

func (c *ClientCLI) StartClient() {

	c.ClearScreen()

	for {
		action := c.promptForMainAction()
		switch action {
		case CreateGameAction:
			c.handleCreateGame()
		case SelectGameAction:
			c.handleSelectGame()
		case ExitAction:
			return
		}
	}
}

func (c *ClientCLI) handleCreateGame() {
	difficulty := c.promptForDifficulty()

	game, gameID := c.Controller().CreateGame(difficulty)

	fmt.Printf("New game created with ID: %s\n", gameID)

	c.currentGame = game
	c.currentGameID = gameID

	c.handleGame()
}

func (c *ClientCLI) handleSelectGame() {
	gameID := c.promptForGameID()

	game, err := c.Controller().GetGame(gameID)

	if err != nil {
		fmt.Println("Game not found!")
		return
	}

	c.currentGame = game
	c.currentGameID = gameID

	c.handleGame()
}

func (c *ClientCLI) handleGame() {

	c.ClearScreen()

	c.printCurrentGame()

	for {
		action := c.promptForGameAction()
		switch action {
		case ShowBoardAction:
			c.ClearScreen()
			c.printCurrentGame()
		case ShowGameHistoryAction:
			c.ClearScreen()
			c.printCurrentGameHistory()
		case MakeMoveAction:
			row, col, value := c.promptForMove()

			if err := c.Controller().MakeMove(c.currentGame, row, col, value); err != nil {
				fmt.Printf("Error: %v\n\n", err)
			} else {
				c.ClearScreen()
				c.printCurrentGame()
				fmt.Printf("Move added.\n\n")
			}
		case RemoveMoveAction:
			row, col := c.promptForCell()
			if err := c.Controller().RemoveMove(c.currentGame, row, col); err != nil {
				fmt.Printf("Error: %v\n\n", err)
			} else {
				c.ClearScreen()
				c.printCurrentGame()
				fmt.Printf("Move removed.\n\n")
			}

		case GetHintAction:
			hint, err := c.Controller().GetHint(c.currentGame)
			if err != nil {
				fmt.Printf("Error: %v\n\n", err)
			} else {
				fmt.Printf("Hint: Set value %d at row %d, column %d.\n\n", hint.Value, hint.Row+1, hint.Col+1)
			}
		case SolveAction:
			if err := c.Controller().SolveGame(c.currentGame); err != nil {
				fmt.Printf("Error: %v\n\n", err)
			} else {
				fmt.Println("Game solved. This is your game: ")
				c.printCurrentBoard()
			}
		case BackAction:
			c.ClearScreen()
			return
		}
	}
}

func (c *ClientCLI) printBoard(board Board) {
	fmt.Println(board.String())
}

func (c *ClientCLI) printCurrentBoard() {
	if c.currentGame == nil {
		fmt.Println("No game selected.")
		return
	}

	c.printBoard((*c.currentGame).Board)
}

func (c *ClientCLI) printGame(game Game) {
	fmt.Println(game.String())
}

func (c *ClientCLI) printCurrentGame() {
	if c.currentGame == nil {
		fmt.Println("No game selected.")
		return
	}

	c.printGame(*c.currentGame)
}

func (c *ClientCLI) printGameHistory(game Game) {

	boardTimeline, err := c.Controller().GetBoardHistory(&game)

	if err != nil {
		fmt.Printf("Error getting game history.\n\n")
		return
	}

	fmt.Printf("Game history for game %s\n\n", game.ID)

	for _, boardHistory := range boardTimeline {

		timestamp := time.Unix(boardHistory.Timestamp, 0).Format(time.Stamp)

		if boardHistory.ActionType == BoardRemovedHistoryAction {
			fmt.Printf("Move %d [%s] - Removed value  %d at row %d, column %d\n\n", boardHistory.Index, timestamp, boardHistory.PreviousValue, boardHistory.Row+1, boardHistory.Col+1)
		} else if boardHistory.ActionType == BoardAddedHistoryAction {
			fmt.Printf("Move %d [%s] - Set value %d at row %d, column %d\n\n", boardHistory.Index, timestamp, boardHistory.CurrentValue, boardHistory.Row+1, boardHistory.Col+1)
		} else if boardHistory.ActionType == BoardCreatedHistoryAction {
			fmt.Printf("Game created [%s]\n\n", timestamp)
		}

		c.printBoard(boardHistory.Board)
	}
}

func (c *ClientCLI) printCurrentGameHistory() {
	if c.currentGame == nil {
		fmt.Println("No game selected.")
		return
	}

	c.printGameHistory(*c.currentGame)
}

func (c *ClientCLI) promptForMainAction() string {

	Log.Debug().Msg("Prompting for main action")

	prompt := promptui.Select{
		Label: "Select an action",
		Items: []string{CreateGameAction, SelectGameAction, ExitAction},
	}
	_, action, _ := prompt.Run()
	return action
}

func (c *ClientCLI) promptForDifficulty() Difficulty {

	Log.Debug().Msg("Prompting for difficulty")

	prompt := promptui.Select{
		Label: "Select Difficulty",
		Items: []string{"Easy", "Medium", "Hard", "Expert", "Master", "Sadistic"},
	}
	_, result, _ := prompt.Run()

	difficulty, err := StringToDifficulty(result)

	if err != nil {
		Log.Error().Err(err).Str("result", result).Msg("Failed to parse difficulty")
		return c.defaultDifficulty
	}

	Log.Debug().Str("result", result).Msg("Selected difficulty")
	return difficulty
}

func (c *ClientCLI) promptForGameAction() string {

	Log.Debug().Msg("Prompting for game action")

	prompt := promptui.Select{
		Label: "Select an action",
		Items: []string{ShowBoardAction, ShowGameHistoryAction, MakeMoveAction, RemoveMoveAction, GetHintAction, SolveAction, BackAction},
	}

	_, action, _ := prompt.Run()

	Log.Debug().Str("action", action).Msg("Selected action")

	return action
}

func (c *ClientCLI) promptForGameID() string {

	Log.Debug().Msg("Prompting for game ID")

	prompt := promptui.Prompt{
		Label: "Enter Game ID",
	}

	gameID, _ := prompt.Run()

	Log.Debug().Str("game_id", gameID).Msg("Entered game ID")

	return gameID
}

func (c *ClientCLI) promptForMove() (int, int, int) {
	row := c.promptForCoordinate("row")
	col := c.promptForCoordinate("column")
	value := c.promptForValue()
	return row, col, value
}

func (c *ClientCLI) promptForCell() (int, int) {
	row := c.promptForCoordinate("row")
	col := c.promptForCoordinate("column")
	return row, col
}

func (c *ClientCLI) promptForCoordinate(label string) int {

	Log.Debug().Str("label", label).Msg("Prompting for coordinate")

	prompt := promptui.Prompt{
		Label: fmt.Sprintf("Enter %s (1-9)", label),
		Validate: func(input string) error {
			num, err := strconv.Atoi(input)
			if err != nil || num < 1 || num > 9 {
				return fmt.Errorf("invalid %s", label)
			}
			return nil
		},
	}
	coord, _ := prompt.Run()
	num, _ := strconv.Atoi(coord)

	Log.Debug().Int("num", num).Str("label", label).Msg("Entered coordinate")

	return num - 1
}

func (c *ClientCLI) promptForValue() int {

	Log.Debug().Msg("Prompting for value")

	prompt := promptui.Prompt{
		Label: "Enter value (1-9)",
		Validate: func(input string) error {
			num, err := strconv.Atoi(input)
			if err != nil || num < 1 || num > 9 {
				return fmt.Errorf("invalid value")
			}
			return nil
		},
	}

	value, _ := prompt.Run()
	num, _ := strconv.Atoi(value)

	Log.Debug().Int("num", num).Msg("Entered value")

	return num
}

func (c *ClientCLI) printLogo() {
	fmt.Println(logo.Get())
}

func (c *ClientCLI) ClearScreen() {
	screen.Clear()
	screen.MoveTopLeft()
	c.printLogo()
}
