package painter

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/exp/shiny/screen"
)

type State struct {
	BackgroundColor color.Color
	BgRect          *image.Rectangle
	Figures         []image.Point
}

func DefaultState() *State {
	return &State{
		BackgroundColor: color.Black,
		Figures:         []image.Point{},
	}
}

type Operation interface {
	Do(t screen.Texture, s *State) (needsUpdate bool)
}

type OperationList []Operation

func (ol OperationList) Do(t screen.Texture, s *State) (needsUpdate bool) {
	for _, o := range ol {
		if o.Do(t, s) {
			needsUpdate = true
		}
	}
	return
}

var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture, s *State) bool {
	return true
}

type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture, s *State) bool {
	f(t)
	return false
}

type WhiteOp struct{}

func (op WhiteOp) Do(t screen.Texture, s *State) bool {
	s.BackgroundColor = color.White
	return false
}

type GreenOp struct{}

func (op GreenOp) Do(t screen.Texture, s *State) bool {
	s.BackgroundColor = color.RGBA{G: 0xff, A: 0xff}
	return false
}

type BgRectOp struct {
	X1, Y1, X2, Y2 float64
}

func (op BgRectOp) Do(t screen.Texture, s *State) bool {
	size := t.Bounds().Size()
	x1 := int(op.X1 * float64(size.X))
	y1 := int(op.Y1 * float64(size.Y))
	x2 := int(op.X2 * float64(size.X))
	y2 := int(op.Y2 * float64(size.Y))

	rect := image.Rect(x1, y1, x2, y2)
	s.BgRect = &rect
	return false
}

type FigureOp struct {
	X, Y float64
}

func (op FigureOp) Do(t screen.Texture, s *State) bool {
	size := t.Bounds().Size()
	x := int(op.X * float64(size.X))
	y := int(op.Y * float64(size.Y))
	s.Figures = append(s.Figures, image.Point{X: x, Y: y})
	return false
}

type MoveOp struct {
	X, Y float64
}

func (op MoveOp) Do(t screen.Texture, s *State) bool {
	size := t.Bounds().Size()
	dx := int(op.X * float64(size.X))
	dy := int(op.Y * float64(size.Y))

	for i := range s.Figures {
		s.Figures[i].X += dx
		s.Figures[i].Y += dy
	}
	return false
}

type ResetOp struct{}

func (op ResetOp) Do(t screen.Texture, s *State) bool {
	s.BackgroundColor = color.Black
	s.BgRect = nil
	s.Figures = []image.Point{}
	return false
}

type DrawStateOp struct{}

func (op DrawStateOp) Do(t screen.Texture, s *State) bool {
	bounds := t.Bounds()
	t.Fill(bounds, s.BackgroundColor, draw.Src)

	if s.BgRect != nil {
		t.Fill(*s.BgRect, color.Black, draw.Src)
	}

	figureColor := color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xff}
	figureSize := 200

	for _, center := range s.Figures {
		halfFigureSize := figureSize / 2
		horizontalRect := image.Rectangle{
			Min: image.Point{X: center.X - halfFigureSize, Y: center.Y - halfFigureSize/3},
			Max: image.Point{X: center.X + halfFigureSize, Y: center.Y + halfFigureSize/3},
		}
		verticalRect := image.Rectangle{
			Min: image.Point{X: center.X - halfFigureSize/3, Y: center.Y - halfFigureSize},
			Max: image.Point{X: center.X + halfFigureSize/3, Y: center.Y + halfFigureSize},
		}
		t.Fill(horizontalRect, figureColor, draw.Src)
		t.Fill(verticalRect, figureColor, draw.Src)
	}

	return false
}
