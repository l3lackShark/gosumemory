package opengl

import (
	"time"

	"github.com/k0kubun/pp"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

//Init opengl renderer
func Init() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		pp.Println("SDL INIT ERROR: ", err)
		return
	}

	window, err := sdl.CreateWindow("Leaderboard v0.x-alpha",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		screenWidth,
		screenHeight,
		sdl.WINDOW_OPENGL)
	if err != nil {
		pp.Println("SDL WINDOW ERROR: ", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		pp.Println("SDL RENDERER ERROR: ", err)
	}
	defer renderer.Destroy()
	for {
		renderer.SetDrawColor(255, 15, 15, 255)
		//Fill the frame with the set color above
		renderer.Clear()

		//Push the frame
		renderer.Present()
		time.Sleep(1 * time.Second)
	}
}
