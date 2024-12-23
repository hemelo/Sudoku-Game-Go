package main

import (
	"Sudoku-Solver/clients/web"
	"Sudoku-Solver/internals"
	"Sudoku-Solver/pkg/logger"
	"flag"
	"time"
)

var Log = logger.Get()

var Debug = flag.Bool("debug", false, "enable debug logging")

func main() {

	flag.Parse()

	if *Debug {
		logger.SetDebug(true)
	}

	// Initialize the cache and game manager

	gameManager := internals.NewGameManager()
	gameController := internals.NewGameController(gameManager)

	client := internals.Client{
		GameController:    gameController,
		GameManager:       gameManager,
		DefaultDifficulty: internals.Easy,
	}

	webClientOpts := web.ClientWebOpts{
		Host:    "localhost",
		Port:    8080,
		Timeout: 30 * time.Second,
	}

	clientWeb := web.NewClientWeb(client, webClientOpts)
	clientWeb.StartClient()
	//clientCLI := cli.NewClientCLI(client)
	//clientCLI.StartClient()
}
