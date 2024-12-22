package main

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

// NewGameManager initializes a new GameManager with the given cache
func NewGameManager() *GameManager {

	gameCache := NewMemoryCache[*Game]()
	originalBoardCache := NewMemoryCache[Board]()

	return &GameManager{gameCache: gameCache, originalBoardCache: originalBoardCache}
}

// CreateGame creates a new game and saves it in the cache
func (gm *GameManager) CreateGame(gameGenOpts GameGeneratorOpts) (*Game, string) {

	if gameGenOpts.BoardId == "" {
		gameGenOpts.BoardId = uuid.NewString()
	}

	gameID := gameGenOpts.BoardId

	Log.Debug().Str("game_id", gameID).Str("difficulty", gameGenOpts.Difficulty.String()).Msg("creating new game")

	board, err := generateBoard(gameGenOpts.BoardGeneratorOpts)

	if err != nil {
		Log.Error().Err(err).Msg("failed to generate board")
		return nil, ""
	}

	game := &Game{
		Board:      board,
		Difficulty: gameGenOpts.Difficulty,
		ID:         gameID,
		CreatedAt:  time.Now().Unix(),
	}

	defer gm.SaveGame(game)

	Log.Info().Str("game_id", gameID).Str("difficulty", gameGenOpts.Difficulty.String()).Msg("game created")

	return game, gameID
}

// GetGame retrieves a game by ID
func (gm *GameManager) GetGame(gameID string) (*Game, error) {

	Log.Debug().Str("game_id", gameID).Msg("getting game from cache")

	game, err := gm.gameCache.Load(gameID)

	if err != nil {
		Log.Error().Str("game_id", gameID).Err(err).Msg("failed to get game from cache")
		return nil, fmt.Errorf("game not found: %v", err)
	}

	return game, nil
}

func (gm *GameManager) SaveGame(game *Game) {

	Log.Debug().Str("game_id", game.ID).Msg("saving game in cache")

	_, err := gm.gameCache.Load(game.ID)

	if err != nil {
		err := gm.originalBoardCache.Save(game.ID, game.Board)

		if err != nil {
			Log.Error().Str("game_id", game.ID).Err(err).Msg("failed to save original board in cache")
		}
	}

	err = gm.gameCache.Save(game.ID, game)

	if err != nil {
		Log.Error().Str("game_id", game.ID).Err(err).Msg("failed to save game in cache")
	}

	Log.Info().Str("game_id", game.ID).Msg("game saved in cache")
}

// DeleteGame deletes a game by ID
func (gm *GameManager) DeleteGame(gameID string) error {

	Log.Debug().Str("game_id", gameID).Msg("deleting game from cache")

	err := gm.gameCache.Delete(gameID)

	if err != nil {
		Log.Error().Str("game_id", gameID).Err(err).Msg("failed to delete game from cache")
		return fmt.Errorf("failed to delete game: %v", err)
	}

	Log.Info().Str("game_id", gameID).Msg("game deleted from cache")

	err = gm.originalBoardCache.Delete(gameID)

	if err != nil {
		Log.Error().Str("game_id", gameID).Err(err).Msg("failed to delete original board from cache")
	}

	return nil
}

func (gm *GameManager) GetOriginalBoard(gameID string) (Board, error) {

	Log.Debug().Str("game_id", gameID).Msg("getting original board from cache")

	board, err := gm.originalBoardCache.Load(gameID)

	if err != nil {
		Log.Error().Str("game_id", gameID).Err(err).Msg("failed to get original board from cache")
		return Board{}, fmt.Errorf("original board not found: %v", err)
	}

	return board, nil
}
