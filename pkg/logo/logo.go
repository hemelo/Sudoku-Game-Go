package logo

import (
	"github.com/mbndr/figlet4go"
	"sync"
)

var once sync.Once

var logo string

func Get() string {
	once.Do(func() {
		render()
	})

	return logo
}

func render() {
	ascii := figlet4go.NewAsciiRender()

	// Adding the colors to RenderOptions
	options := figlet4go.NewRenderOptions()

	options.FontColor = []figlet4go.Color{
		// Colors can be given by default ansi color codes...
		figlet4go.ColorGreen,
		figlet4go.ColorYellow,
		figlet4go.ColorCyan,
	}

	renderStr, _ := ascii.RenderOpts("Sudoku", options)
	logo = renderStr
}
