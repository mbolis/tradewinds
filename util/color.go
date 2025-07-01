package util

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func Col(c color.Color) *ebiten.Image {
	i := ebiten.NewImage(3, 3)
	i.Fill(c)
	return i.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
}
