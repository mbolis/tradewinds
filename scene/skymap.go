package scene

import (
	"bytes"
	"image/color"
	"log"
	"slices"

	"github.com/chewxy/math32"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mbolis/tradewinds/model"
	"github.com/mbolis/tradewinds/util"
)

type Skymap struct {
	// TODO mode: inspect, plan, travel
	// TODO graphics
	// TODO viewport: zoom, pan
	// TODO POIs: cities/towns, ports, facilities
	selectedPOI *POI
	Player      *model.Player
}

var _ Scene = (*Skymap)(nil)

type POI struct {
	Name  string
	Color color.Color
	X, Y  float32

	fillVertices, strokeVertices []ebiten.Vertex
	fillIndices, strokeIndices   []uint16
}

func (poi *POI) Draw(screen *ebiten.Image) {
	if poi.fillVertices == nil {
		poi.fillVertices = slices.Clone(poiFillVertices)
		poi.fillIndices = slices.Clone(poiFillIndices)
		poi.strokeVertices = slices.Clone(poiStrokeVertices)
		poi.strokeIndices = slices.Clone(poiStrokeIndices)

		for i := range poi.fillVertices {
			poi.fillVertices[i].DstX += float32(poi.X)
			poi.fillVertices[i].DstY += float32(poi.Y)

			r, g, b, _ := poi.Color.RGBA()
			poi.fillVertices[i].ColorR = float32(r) / float32(0xffff)
			poi.fillVertices[i].ColorG = float32(g) / float32(0xffff)
			poi.fillVertices[i].ColorB = float32(b) / float32(0xffff)
			poi.fillVertices[i].ColorA = 1
		}

		for i := range poi.strokeVertices {
			poi.strokeVertices[i].DstX += float32(poi.X)
			poi.strokeVertices[i].DstY += float32(poi.Y)
		}
	}

	screen.DrawTriangles(poi.fillVertices, poi.fillIndices, whiteImg, &ebiten.DrawTrianglesOptions{FillRule: ebiten.FillRuleNonZero, AntiAlias: true})
	screen.DrawTriangles(poi.strokeVertices, poi.strokeIndices, whiteImg, &ebiten.DrawTrianglesOptions{FillRule: ebiten.FillRuleFillAll, AntiAlias: true})

	var opt text.DrawOptions
	opt.PrimaryAlign = text.AlignCenter
	opt.SecondaryAlign = text.AlignCenter
	opt.GeoM.Translate(float64(poi.X)+1, float64(poi.Y)+11)
	opt.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, poi.Name, font, &opt)

	opt.GeoM.Translate(-1, -1)
	opt.ColorScale = ebiten.ColorScale{}
	opt.ColorScale.ScaleWithColor(poi.Color)
	text.Draw(screen, poi.Name, font, &opt)
}

var (
	skyBlue    = color.NRGBA{R: 103, G: 202, B: 255, A: 255}
	grassGreen = color.NRGBA{R: 118, G: 220, B: 0, A: 255}
	sandBrown  = color.NRGBA{R: 237, G: 208, B: 36, A: 255}

	islands = ebiten.NewImage(640, 480)

	whiteImg  = util.Col(color.White)
	darkRed   = color.NRGBA{R: 127, G: 0, B: 0, A: 255}
	red       = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	orange    = color.NRGBA{R: 255, G: 127, B: 0, A: 255}
	yellow    = color.NRGBA{R: 255, G: 191, B: 0, A: 255}
	lime      = color.NRGBA{R: 191, G: 255, B: 0, A: 255}
	green     = color.NRGBA{R: 0, G: 191, B: 0, A: 255}
	teal      = color.NRGBA{R: 0, G: 191, B: 255, A: 255}
	blue      = color.NRGBA{R: 31, G: 31, B: 255, A: 255}
	lightBlue = color.NRGBA{R: 63, G: 127, B: 255, A: 255}
	purple    = color.NRGBA{R: 127, G: 0, B: 255, A: 255}
	magenta   = color.NRGBA{R: 255, G: 0, B: 255, A: 255}

	pois = []POI{
		{Name: "Kythera", Color: red, X: 120, Y: 200},
		{Name: "Ogigia", Color: lime, X: 340, Y: 240},
		{Name: "Eritia", Color: blue, X: 600, Y: 180},
		{Name: "Ithaca", Color: yellow, X: 280, Y: 400},
		{Name: "Massalia", Color: purple, X: 230, Y: 60},
		{Name: "Sybaris", Color: magenta, X: 450, Y: 65},
		{Name: "Zancle", Color: darkRed, X: 380, Y: 130},
		{Name: "Acragas", Color: teal, X: 180, Y: 440},
		{Name: "Himaera", Color: orange, X: 240, Y: 290},
		{Name: "Elea", Color: green, X: 540, Y: 300},
		{Name: "Emporion", Color: lightBlue, X: 600, Y: 450},
	}
	poiFillVertices, poiStrokeVertices []ebiten.Vertex
	poiFillIndices, poiStrokeIndices   []uint16

	fontSource *text.GoTextFaceSource
	font       *text.GoTextFace
)

func init() {
	// islands
	grassland := ebiten.NewImage(640, 480)
	for _, island := range []struct{ x, y, r float32 }{
		{100, 400, 100},
		{220, 340, 60},
		{280, 400, 10},
		{600, 440, 160},
		{340, 30, 120},
		{180, 10, 80},
		{520, 40, 80},
		{600, 180, 12},
		{340, 240, 8},
		{120, 200, 16},
	} {
		vector.DrawFilledCircle(islands, island.x, island.y, island.r+2, sandBrown, true)
		vector.DrawFilledCircle(grassland, island.x, island.y, island.r-2, grassGreen, true)
	}
	islands.DrawImage(grassland, nil)

	// POI sprite buffer
	var path vector.Path
	path.MoveTo(-4, +2)
	path.LineTo(-4, -4)
	path.LineTo(-8, -4)
	path.LineTo(0, -10)
	path.LineTo(+8, -4)
	path.LineTo(+4, -4)
	path.LineTo(+4, +2)
	path.Close()

	poiFillVertices, poiFillIndices = path.AppendVerticesAndIndicesForFilling(poiFillVertices, poiFillIndices)
	poiStrokeVertices, poiStrokeIndices = path.AppendVerticesAndIndicesForStroke(poiStrokeVertices, poiStrokeIndices, &vector.StrokeOptions{
		Width: 1,
	})

	for i := range poiFillVertices {
		poiFillVertices[i].SrcX = 1
		poiFillVertices[i].SrcY = 1
	}
	for i := range poiStrokeVertices {
		poiStrokeVertices[i].SrcX = 1
		poiStrokeVertices[i].SrcY = 1

		poiStrokeVertices[i].ColorR = 0
		poiStrokeVertices[i].ColorG = 0
		poiStrokeVertices[i].ColorB = 0
		poiStrokeVertices[i].ColorA = 1
	}

	// fonts
	var err error
	fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	font = &text.GoTextFace{
		Source: fontSource,
		Size:   14,
	}
}

func (s *Skymap) Draw(screen *ebiten.Image) {
	screen.Fill(skyBlue)
	screen.DrawImage(islands, nil)
	for _, poi := range pois {
		poi.Draw(screen)
	}
	s.Player.Draw(screen)

	if s.selectedPOI != nil {
		drawArrow(screen, s.Player.X, s.Player.Y, s.selectedPOI.X, s.selectedPOI.Y)
	}
}

func drawArrow(image *ebiten.Image, x0, y0, x1, y1 float32) {
	dx := x1 - x0
	dy := y1 - y0
	distance := math32.Sqrt(dx*dx + dy*dy)
	if distance == 0 {
		return
	}

	angleLine := -math32.Atan2(dy, dx)

	xOffset := 6 * math32.Cos(angleLine)
	yOffset := 6 * math32.Sin(angleLine)

	var path vector.Path
	path.MoveTo(x0+xOffset, y0-yOffset)
	path.LineTo(x1-xOffset, y1+yOffset)
	// TODO dashed rainbow line

	angleHeadLeft := angleLine + 0.75*math32.Pi
	angleHeadRight := angleLine - 0.75*math32.Pi
	path.LineTo(x1-xOffset+10*math32.Cos(angleHeadLeft), y1+yOffset-10*math32.Sin(angleHeadLeft))
	path.MoveTo(x1-xOffset, y1+yOffset)
	path.LineTo(x1-xOffset+10*math32.Cos(angleHeadRight), y1+yOffset-10*math32.Sin(angleHeadRight))

	vertices, indices := path.AppendVerticesAndIndicesForStroke([]ebiten.Vertex{}, []uint16{}, &vector.StrokeOptions{
		Width:   4,
		LineCap: vector.LineCapRound,
	})
	verticesShadow, indicesShadow := path.AppendVerticesAndIndicesForStroke([]ebiten.Vertex{}, []uint16{}, &vector.StrokeOptions{
		Width:   5,
		LineCap: vector.LineCapRound,
	})

	var c color.Color
	switch {
	case distance <= 100:
		c = color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	case distance <= 300:
		c = color.NRGBA{R: uint8(255 * (distance - 100) / 200), G: 255, B: 0, A: 255}
	default:
		c = color.NRGBA{R: 255, G: 255 - uint8(255*min((distance-300)/200, 1)), B: 0, A: 255}
	}
	r, g, b, a := c.RGBA()
	for i := range vertices {
		vertices[i].ColorR = float32(r) / 0xffff
		vertices[i].ColorG = float32(g) / 0xffff
		vertices[i].ColorB = float32(b) / 0xffff
		vertices[i].ColorA = float32(a) / 0xffff
	}
	for i := range verticesShadow {
		verticesShadow[i].ColorR = 0
		verticesShadow[i].ColorG = 0
		verticesShadow[i].ColorB = 0
		verticesShadow[i].ColorA = 0.3
	}

	image.DrawTriangles(verticesShadow, indicesShadow, whiteImg, &ebiten.DrawTrianglesOptions{FillRule: ebiten.FillRuleFillAll, AntiAlias: true})
	image.DrawTriangles(vertices, indices, whiteImg, &ebiten.DrawTrianglesOptions{FillRule: ebiten.FillRuleFillAll, AntiAlias: true})
}

func (s *Skymap) Update() error {
	cursorShape := ebiten.CursorShape()
	hoveredPOI := getHoveredPOI()
	if hoveredPOI != nil {
		if cursorShape != ebiten.CursorShapePointer {
			ebiten.SetCursorShape(ebiten.CursorShapePointer)
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			s.selectedPOI = hoveredPOI
		}
	} else {
		if cursorShape != ebiten.CursorShapeDefault {
			ebiten.SetCursorShape(ebiten.CursorShapeDefault)
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			s.selectedPOI = nil
		}
	}
	return nil
}

func getHoveredPOI() *POI {
	const hoverRange2 = 100
	mouseX, mouseY := ebiten.CursorPosition()

	for _, poi := range pois {
		dx := float32(mouseX) - poi.X
		dy := float32(mouseY) - poi.Y
		if dx*dx+dy*dy <= hoverRange2 {
			return &poi
		}
	}
	return nil
}
