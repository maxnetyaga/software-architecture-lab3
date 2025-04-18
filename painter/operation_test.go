package painter_test

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"

	"golang.org/x/exp/shiny/screen"

	"github.com/maxnetyaga/software-architecture-lab3/painter"
)

type fillCall struct {
	dr image.Rectangle
	src color.Color
	op draw.Op
}

type mockTexture struct {
	size   image.Point
	buffer *image.RGBA
	fillCalls []fillCall
}

var testTextureSize = image.Pt(800, 800)

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return m.size }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.size}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}

func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	draw.Draw(m.buffer, dr, image.NewUniform(src), image.Point{}, op)
	m.fillCalls = append(m.fillCalls, fillCall{dr, src, op})
}

func newMockTexture(size image.Point) *mockTexture {
	return &mockTexture{
		size: size,
		buffer: image.NewRGBA(image.Rectangle{Max: size}),
		fillCalls: []fillCall{},
	}
}

func checkPixelColor(t *testing.T, texture *mockTexture, x, y int, expected color.Color, message string) {
	gotR, gotG, gotB, gotA := texture.buffer.At(x, y).RGBA()
	expectedR, expectedG, expectedB, expectedA := expected.RGBA()

	if gotR != expectedR || gotG != expectedG || gotB != expectedB || gotA != expectedA {
		t.Errorf("%s at (%d, %d). Expected RGBA: %v, Got RGBA: %v (Original Got: %v)",
			message, x, y,
			color.RGBA64{uint16(expectedR), uint16(expectedG), uint16(expectedB), uint16(expectedA)},
			color.RGBA64{uint16(gotR), uint16(gotG), uint16(gotB), uint16(gotA)},
			texture.buffer.At(x, y),
		)
	}
}

func checkState(t *testing.T, state *painter.State, expected painter.State, message string) {
	if state.BackgroundColor != expected.BackgroundColor {
		t.Errorf("%s: BackgroundColor mismatch. Expected: %v, Got: %v", message, expected.BackgroundColor, state.BackgroundColor)
	}
	if (state.BgRect == nil && expected.BgRect != nil) || (state.BgRect != nil && expected.BgRect == nil) {
		 t.Errorf("%s: BgRect mismatch. Expected: %v, Got: %v", message, expected.BgRect, state.BgRect)
	} else if state.BgRect != nil && expected.BgRect != nil {
		if *state.BgRect != *expected.BgRect {
			 t.Errorf("%s: BgRect value mismatch. Expected: %v, Got: %v", message, *expected.BgRect, *state.BgRect)
		}
	}
	if !reflect.DeepEqual(state.Figures, expected.Figures) {
		t.Errorf("%s: Figures mismatch. Expected: %v, Got: %v", message, expected.Figures, state.Figures)
	}
}


func TestOperations_Do(t *testing.T) {

	t.Run("WhiteOp", func(t *testing.T) {
		state := painter.DefaultState()
		texture := newMockTexture(testTextureSize)

		op := painter.WhiteOp{}
		needsUpdate := op.Do(texture, state)

		checkState(t, state, painter.State{BackgroundColor: color.White, BgRect: nil, Figures: []image.Point{}}, "State after WhiteOp")
		if needsUpdate {
			t.Error("WhiteOp returned needsUpdate = true unexpectedly")
		}
		if len(texture.fillCalls) != 0 {
			t.Errorf("WhiteOp called Fill %d times, expected 0", len(texture.fillCalls))
		}
	})

	t.Run("GreenOp", func(t *testing.T) {
		state := painter.DefaultState()
		texture := newMockTexture(testTextureSize)

		op := painter.GreenOp{}
		needsUpdate := op.Do(texture, state)

		greenColor := color.RGBA{G: 255, A: 255}
		checkState(t, state, painter.State{BackgroundColor: greenColor, BgRect: nil, Figures: []image.Point{}}, "State after GreenOp")
		if needsUpdate {
			t.Error("GreenOp returned needsUpdate = true unexpectedly")
		}
		if len(texture.fillCalls) != 0 {
			t.Errorf("GreenOp called Fill %d times, expected 0", len(texture.fillCalls))
		}
	})

	t.Run("BgRectOp", func(t *testing.T) {
		state := painter.DefaultState()
		texture := newMockTexture(testTextureSize)

		op := painter.BgRectOp{X1: 0.1, Y1: 0.1, X2: 0.9, Y2: 0.9}
		needsUpdate := op.Do(texture, state)

		expectedRect := image.Rect(int(0.1*800), int(0.1*800), int(0.9*800), int(0.9*800))
		checkState(t, state, painter.State{BackgroundColor: color.Black, BgRect: &expectedRect, Figures: []image.Point{}}, "State after BgRectOp")
		if needsUpdate {
			t.Error("BgRectOp returned needsUpdate = true unexpectedly")
		}
		if len(texture.fillCalls) != 0 {
			t.Errorf("BgRectOp called Fill %d times, expected 0", len(texture.fillCalls))
		}

		op2 := painter.BgRectOp{X1: 0.2, Y1: 0.2, X2: 0.8, Y2: 0.8}
		op2.Do(texture, state)
		expectedRect2 := image.Rect(int(0.2*800), int(0.2*800), int(0.8*800), int(0.8*800))
		checkState(t, state, painter.State{BackgroundColor: color.Black, BgRect: &expectedRect2, Figures: []image.Point{}}, "State after second BgRectOp")
	})

	t.Run("FigureOp", func(t *testing.T) {
		state := painter.DefaultState()
		texture := newMockTexture(testTextureSize)

		op1 := painter.FigureOp{X: 0.5, Y: 0.5}
		needsUpdate1 := op1.Do(texture, state)

		expectedFigures1 := []image.Point{{X: int(0.5*800), Y: int(0.5*800)}}
		checkState(t, state, painter.State{BackgroundColor: color.Black, BgRect: nil, Figures: expectedFigures1}, "State after first FigureOp")
		if needsUpdate1 {
			t.Error("FigureOp returned needsUpdate = true unexpectedly")
		}
		if len(texture.fillCalls) != 0 {
			t.Errorf("FigureOp called Fill %d times, expected 0", len(texture.fillCalls))
		}

		op2 := painter.FigureOp{X: 0.2, Y: 0.8}
		needsUpdate2 := op2.Do(texture, state)

		expectedFigures2 := []image.Point{{X: int(0.5*800), Y: int(0.5*800)}, {X: int(0.2*800), Y: int(0.8*800)}}
		checkState(t, state, painter.State{BackgroundColor: color.Black, BgRect: nil, Figures: expectedFigures2}, "State after second FigureOp")
		if needsUpdate2 {
			t.Error("Second FigureOp returned needsUpdate = true unexpectedly")
		}
		 if len(texture.fillCalls) != 0 {
			t.Errorf("FigureOp called Fill %d times, expected 0", len(texture.fillCalls))
		}
	})

	t.Run("MoveOp", func(t *testing.T) {
		state := painter.DefaultState()
		state.Figures = []image.Point{{X: 400, Y: 400}, {X: 100, Y: 100}}
		texture := newMockTexture(testTextureSize)

		op := painter.MoveOp{X: 0.1, Y: 0.1}
		needsUpdate := op.Do(texture, state)

		expectedFigures := []image.Point{{X: 480, Y: 480}, {X: 180, Y: 180}}
		checkState(t, state, painter.State{BackgroundColor: color.Black, BgRect: nil, Figures: expectedFigures}, "State after MoveOp")
		if needsUpdate {
			t.Error("MoveOp returned needsUpdate = true unexpectedly")
		}
		if len(texture.fillCalls) != 0 {
			t.Errorf("MoveOp called Fill %d times, expected 0", len(texture.fillCalls))
		}
	})

	t.Run("ResetOp", func(t *testing.T) {
		state := &painter.State{
			BackgroundColor: color.RGBA{G: 255, A: 255},
			BgRect:          &image.Rectangle{Min: image.Point{100, 100}, Max: image.Point{700, 700}},
			Figures:         []image.Point{{200, 200}, {600, 600}},
		}
		texture := newMockTexture(testTextureSize)

		op := painter.ResetOp{}
		needsUpdate := op.Do(texture, state)

		checkState(t, state, painter.State{BackgroundColor: color.Black, BgRect: nil, Figures: []image.Point{}}, "State after ResetOp")
		if needsUpdate {
			t.Error("ResetOp returned needsUpdate = true unexpectedly")
		}
		if len(texture.fillCalls) != 0 {
			t.Errorf("ResetOp called Fill %d times, expected 0", len(texture.fillCalls))
		}
	})

	t.Run("UpdateOp", func(t *testing.T) {
		state := painter.DefaultState()
		texture := newMockTexture(testTextureSize)

		op := painter.UpdateOp
		needsUpdate := op.Do(texture, state)

		checkState(t, state, painter.State{BackgroundColor: color.Black, BgRect: nil, Figures: []image.Point{}}, "State after UpdateOp")
		if len(texture.fillCalls) != 0 {
			t.Errorf("UpdateOp called Fill %d times, expected 0", len(texture.fillCalls))
		}

		if !needsUpdate {
			t.Error("UpdateOp returned needsUpdate = false unexpectedly")
		}
	})

	t.Run("DrawStateOp", func(t *testing.T) {
		texture := newMockTexture(testTextureSize)
		state := &painter.State{
			BackgroundColor: color.RGBA{G: 255, A: 255},
			BgRect:          &image.Rectangle{Min: image.Point{200, 200}, Max: image.Point{600, 600}},
			Figures:         []image.Point{{400, 400}, {100, 100}},
		}

		op := painter.DrawStateOp{}
		needsUpdate := op.Do(texture, state)

		checkState(t, state, *state, "State after DrawStateOp")
		 if needsUpdate {
			t.Error("DrawStateOp returned needsUpdate = true unexpectedly")
		}

		checkPixelColor(t, texture, 10, 10, color.RGBA{G: 255, A: 255}, "Pixel color mismatch in DrawStateOp (Background)")
		checkPixelColor(t, texture, 300, 300, color.Black, "Pixel color mismatch in DrawStateOp (BgRect)")
		checkPixelColor(t, texture, 400, 400, color.RGBA{R: 255, G: 255, B: 0, A: 255}, "Pixel color mismatch in DrawStateOp (Figure 1)")
		checkPixelColor(t, texture, 100, 100, color.RGBA{R: 255, G: 255, B: 0, A: 255}, "Pixel color mismatch in DrawStateOp (Figure 2)")

		expectedFillCalls := 1 + 1 + (2 * 2)
		if len(texture.fillCalls) != expectedFillCalls {
			 t.Errorf("DrawStateOp called Fill %d times, expected %d", len(texture.fillCalls), expectedFillCalls)
		}
	})

	 t.Run("OperationList", func(t *testing.T) {
		state := painter.DefaultState()
		texture := newMockTexture(testTextureSize)

		opList := painter.OperationList{
			painter.WhiteOp{},
			painter.BgRectOp{X1: 0.3, Y1: 0.3, X2: 0.7, Y2: 0.7},
			painter.FigureOp{X: 0.5, Y: 0.5},
			painter.UpdateOp,
		}

		needsUpdate := opList.Do(texture, state)

		expectedRect := image.Rect(int(0.3*800), int(0.3*800), int(0.7*800), int(0.7*800))
		expectedFigures := []image.Point{{X: int(0.5*800), Y: int(0.5*800)}}
		checkState(t, state, painter.State{BackgroundColor: color.White, BgRect: &expectedRect, Figures: expectedFigures}, "State after OperationList")

		if !needsUpdate {
			t.Error("OperationList returned needsUpdate = false unexpectedly")
		}

		if len(texture.fillCalls) != 0 {
			t.Errorf("OperationList called Fill %d times, expected 0 (Fill calls should happen in Loop with DrawStateOp)", len(texture.fillCalls))
		}
	})

	 t.Run("OperationFunc", func(t *testing.T) {
		state := painter.DefaultState()
		texture := newMockTexture(testTextureSize)
		fillCount := 0

		opFunc := painter.OperationFunc(func(t screen.Texture){
			t.Fill(image.Rect(0,0,1,1), color.RGBA{R: 255, A: 255}, draw.Src)
			fillCount++
		})

		needsUpdate := opFunc.Do(texture, state)

		checkState(t, state, painter.State{BackgroundColor: color.Black, BgRect: nil, Figures: []image.Point{}}, "State after OperationFunc")
		if needsUpdate {
			t.Error("OperationFunc returned needsUpdate = true unexpectedly")
		}
		if fillCount != 1 {
			t.Errorf("OperationFunc internal function called %d times, expected 1", fillCount)
		}
		 if len(texture.fillCalls) != 1 {
			 t.Errorf("OperationFunc called Fill %d times, expected 1", len(texture.fillCalls))
		 } else {
			 call := texture.fillCalls[0]
			 expectedRect := image.Rect(0,0,1,1)
			 expectedRed := color.RGBA{R: 255, A: 255}
			 if call.dr != expectedRect || call.src != expectedRed || call.op != draw.Src {
				 t.Errorf("Fill call parameters mismatch in OperationFunc. Expected: %v, %v, %v, Got: %v, %v, %v",
					 expectedRect, expectedRed, draw.Src,
					 call.dr, call.src, call.op)
			 }
		 }
	})
}