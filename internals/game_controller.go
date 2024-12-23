package internals

import (
	"math/rand"
	"sort"
	"time"
)

func NewGameController(gameManager *GameManager) *GameController {
	return &GameController{gameManager: gameManager}
}

func (gc *GameController) GetGame(gameID string) (*Game, error) {
	return gc.gameManager.GetGame(gameID)
}

func (gc *GameController) MakeMove(game *Game, row, col, value int) error {

	if err := game.MakeMove(row, col, value); err != nil {
		return err
	}

	defer gc.gameManager.SaveGame(game)
	return nil
}

func (gc *GameController) RemoveMove(game *Game, row, col int) error {

	if err := game.RemoveMove(row, col); err != nil {
		return err
	}

	defer gc.gameManager.SaveGame(game)
	return nil
}

func (gc *GameController) SolveGame(game *Game) error {

	if err := game.SolveGame(); err != nil {
		return err
	}

	defer gc.gameManager.SaveGame(game)
	return nil
}

func (gc *GameController) GetHint(game *Game) (*Hint, error) {
	return game.GetHint()
}

func (gc *GameController) CreateGame(difficulty Difficulty) (*Game, string) {

	createOpts := GameGeneratorOpts{
		BoardGeneratorOpts: BoardGeneratorOpts{
			Difficulty: difficulty,
			SeedSourceGenerate: func() *rand.Rand {
				source := rand.NewSource(time.Now().UnixNano())
				return rand.New(source)
			},
		},
	}

	Log.Debug().Msg("Calling manager to create game")

	return gc.gameManager.CreateGame(createOpts)
}

func (gc *GameController) GetBoardHistory(game *Game) (BoardTimeline, error) {

	Log.Info().Str("game_id", game.ID).Msg("Getting board history")

	boards := make(BoardTimeline, 0, len(game.MoveHistory)+1)

	board, err := gc.gameManager.GetOriginalBoard(game.ID)

	if err != nil {
		Log.Error().Err(err).Msg("Failed to get original board")
		return nil, err
	}

	originalBoardHistory := BoardHistory{
		Index:      1,
		Board:      board,
		Timestamp:  game.CreatedAt,
		ActionType: BoardCreatedHistoryAction,
	}

	boards = append(boards, originalBoardHistory)

	history := game.MoveHistory
	sort.Sort(history)

	var actionType BoardHistoryAction

	for index, move := range history {
		board[move.Row][move.Col].Value = move.CurrentValue

		if move.CurrentValue == 0 {
			actionType = BoardRemovedHistoryAction
		} else {
			actionType = BoardAddedHistoryAction
		}

		boardHistory := BoardHistory{
			Index:      index + 1,
			Move:       move,
			Board:      board,
			Timestamp:  move.Timestamp,
			ActionType: actionType,
		}

		boards = append(boards, boardHistory)
	}

	return boards, nil
}
