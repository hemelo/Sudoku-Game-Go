package main

import (
	"Sudoku-Solver/pkg/logger"
	"flag"
)

var Log = logger.Get()

var Debug = flag.Bool("debug", false, "enable debug logging")

func main() {

	flag.Parse()

	if *Debug {
		logger.SetDebug(true)
	}

	// Initialize the cache and game manager

	gameManager := NewGameManager()
	gameController := NewGameController(gameManager)

	clientOpts := ClientCLIOpts{
		gameController:    gameController,
		gameManager:       gameManager,
		defaultDifficulty: Easy,
	}

	clientCLI := NewClientCLI(clientOpts)
	clientCLI.StartClient()
}
