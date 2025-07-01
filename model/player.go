package model

import (
	"image/color"
	"sync"
	"time"

	"github.com/chewxy/math32"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mbolis/tradewinds/util"
)

type Player struct {
	X, Y float32

	fillVertices, strokeVertices []ebiten.Vertex
	fillIndices, strokeIndices   []uint16

	AnimationStart time.Time
}

var (
	initBuffersOnce sync.Once

	playerFillVertices, playerStrokeVertices []ebiten.Vertex
	playerFillIndices, playerStrokeIndices   []uint16

	whiteImg = util.Col(color.White)
)

func init() {
	// player sprite buffer
	var path vector.Path

	path.MoveTo(-4, -4)
	path.CubicTo(-12, -6, -10, -13, 0, -13)
	path.CubicTo(+10, -13, +12, -6, +4, -4)
	path.Close()

	path.MoveTo(-3, -1)
	path.LineTo(-6, -4)
	path.LineTo(+6, -4)
	path.LineTo(+3, -1)
	path.Close()

	playerFillVertices, playerFillIndices = path.AppendVerticesAndIndicesForFilling(playerFillVertices, playerFillIndices)
	playerStrokeVertices, playerStrokeIndices = path.AppendVerticesAndIndicesForStroke(playerStrokeVertices, playerStrokeIndices, &vector.StrokeOptions{
		Width: 1,
	})

	for i := range playerFillVertices {
		playerFillVertices[i].SrcX = 1
		playerFillVertices[i].SrcY = 1

		playerStrokeVertices[i].ColorR = 1
		playerStrokeVertices[i].ColorG = 1
		playerStrokeVertices[i].ColorB = 1
		playerStrokeVertices[i].ColorA = 1
	}
	for i := range playerStrokeVertices {
		playerStrokeVertices[i].SrcX = 1
		playerStrokeVertices[i].SrcY = 1

		playerStrokeVertices[i].ColorR = 0
		playerStrokeVertices[i].ColorG = 0
		playerStrokeVertices[i].ColorB = 0
		playerStrokeVertices[i].ColorA = 1
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	initBuffersOnce.Do(func() {
		p.fillVertices = make([]ebiten.Vertex, len(playerFillVertices))
		p.fillIndices = make([]uint16, len(playerFillIndices))
		p.strokeVertices = make([]ebiten.Vertex, len(playerStrokeVertices))
		p.strokeIndices = make([]uint16, len(playerStrokeIndices))

		copy(p.fillVertices, playerFillVertices)
		copy(p.fillIndices, playerFillIndices)
		copy(p.strokeVertices, playerStrokeVertices)
		copy(p.strokeIndices, playerStrokeIndices)
	})

	deltaT := float32(time.Since(p.AnimationStart).Seconds())
	deltaX := 1.5 * math32.Sin(deltaT/2)
	deltaY := -0.5 * math32.Cos(deltaT/3)

	for i := range p.fillVertices {
		p.fillVertices[i].DstX = playerFillVertices[i].DstX + p.X + deltaX
		p.fillVertices[i].DstY = playerFillVertices[i].DstY + p.Y + deltaY
	}
	for i := range p.strokeVertices {
		p.strokeVertices[i].DstX = playerStrokeVertices[i].DstX + p.X + deltaX
		p.strokeVertices[i].DstY = playerStrokeVertices[i].DstY + p.Y + deltaY
	}

	screen.DrawTriangles(p.fillVertices, p.fillIndices, whiteImg, &ebiten.DrawTrianglesOptions{FillRule: ebiten.FillRuleNonZero, AntiAlias: true})
	screen.DrawTriangles(p.strokeVertices, p.strokeIndices, whiteImg, &ebiten.DrawTrianglesOptions{FillRule: ebiten.FillRuleFillAll, AntiAlias: true})
}
