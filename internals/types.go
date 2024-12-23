package internals

import (
	"Sudoku-Solver/pkg/logger"
	"fmt"
	"math/rand"
	"sync"
)

var Log = logger.Get()

const (
	BoardCreatedHistoryAction = "Board Created"
	BoardAddedHistoryAction   = "Move Added"
	BoardRemovedHistoryAction = "Move Removed"
)

const (
	Easy Difficulty = iota
	Medium
	Hard
	Expert
	Master
	Sadistic
)

type Difficulty int

type SeedSourceGenerator func() *rand.Rand

type Cell struct {
	Value int
	Fixed bool
}

type Board [9][9]Cell

type Move struct {
	Row           int
	Col           int
	PreviousValue int
	CurrentValue  int
	Timestamp     int64
}

type MoveHistory []Move

type Game struct {
	ID          string
	Board       Board
	Difficulty  Difficulty
	MoveHistory MoveHistory
	Active      bool
	Completed   bool
	CreatedAt   int64
}

type Hint struct {
	Row   int
	Col   int
	Value int
}

type BoardGeneratorOpts struct {
	BoardId            string
	Difficulty         Difficulty
	SeedSourceGenerate SeedSourceGenerator
}

type GameGeneratorOpts struct {
	BoardGeneratorOpts
}

type GameManager struct {
	gameCache          Cache[*Game]
	originalBoardCache Cache[Board]
}

type GameController struct {
	gameManager *GameManager
}

type Cache[T any] interface {
	Save(key string, data T) error
	Load(key string) (T, error)
	Delete(key string) error
}

type MemoryCache[T any] struct {
	data map[string]T
	mu   sync.RWMutex
}

type BoardHistoryAction string

type BoardHistory struct {
	Index      int
	Board      Board
	Timestamp  int64
	ActionType BoardHistoryAction
	Move
}

type BoardTimeline []BoardHistory

type Client struct {
	GameController    *GameController
	GameManager       *GameManager
	DefaultDifficulty Difficulty
}

func (a BoardTimeline) Len() int {
	return len(a)
}

func (a BoardTimeline) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a BoardTimeline) Less(i, j int) bool {
	return a[i].Timestamp < a[j].Timestamp
}

func (d Difficulty) String() string {
	switch d {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	case Expert:
		return "Expert"
	case Master:
		return "Master"
	case Sadistic:
		return "Sadistic"
	default:
		return ""
	}
}

func StringToDifficulty(s string) (Difficulty, error) {
	switch s {
	case "Easy":
		return Easy, nil
	case "Medium":
		return Medium, nil
	case "Hard":
		return Hard, nil
	case "Expert":
		return Expert, nil
	case "Master":
		return Master, nil
	case "Sadistic":
		return Sadistic, nil
	default:
		return -1, fmt.Errorf("invalid difficulty: %s", s)
	}
}

func (a MoveHistory) Len() int {
	return len(a)
}

func (a MoveHistory) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a MoveHistory) Less(i, j int) bool {
	return a[i].Timestamp < a[j].Timestamp
}
