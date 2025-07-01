package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mbolis/tradewinds/model"
	"github.com/mbolis/tradewinds/scene"
)

func main() {
	game := &Game{
		scene: &scene.Skymap{
			Player: &model.Player{
				X: 120, Y: 200,
				AnimationStart: time.Now(),
			},
		},
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	scene scene.Scene
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	return g.scene.Update()
}
