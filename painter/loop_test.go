package painter

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
	"time"

	"golang.org/x/exp/shiny/screen"
)

type testReceiver struct {
	lastTexture screen.Texture
	updated     chan struct{}
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
	select {
	case tr.updated <- struct{}{}:
	default:
	}
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("implement me")
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return &mockTexture{size: size, buffer: image.NewRGBA(image.Rectangle{Max: size})}, nil
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("implement me")
}

type mockTexture struct {
	size   image.Point
	buffer *image.RGBA
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return m.size }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.size}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}

func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	draw.Draw(m.buffer, dr, image.NewUniform(src), image.Point{}, op)
}

func checkPixelColor(t *testing.T, texture screen.Texture, x, y int, expected color.Color, message string) {
	mt, ok := texture.(*mockTexture)
	if !ok {
		t.Errorf("Expected *mockTexture, got %T", texture)
		return
	}

	gotR, gotG, gotB, gotA := mt.buffer.At(x, y).RGBA()
	expectedR, expectedG, expectedB, expectedA := expected.RGBA()

	if gotR != expectedR || gotG != expectedG || gotB != expectedB || gotA != expectedA {
		t.Errorf("%s at (%d, %d). Expected RGBA: %v, Got RGBA: %v (Original Got: %v)",
			message, x, y,
			color.RGBA64{R: uint16(expectedR), G: uint16(expectedG), B: uint16(expectedB), A: uint16(expectedA)},
			color.RGBA64{R: uint16(gotR), G: uint16(gotG), B: uint16(gotB), A: uint16(gotA)},
			mt.buffer.At(x, y),
		)
	}
}

func TestLoop_PostAndRun(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	tr.updated = make(chan struct{}, 1)
	l.Receiver = &tr

	mockSc := mockScreen{}
	l.Start(mockSc)
	defer l.StopAndWait()

	t.Run("WhiteBackground", func(t *testing.T) {
		l.Post(ResetOp{})
		l.Post(WhiteOp{})
		l.Post(UpdateOp)

		select {
		case <-tr.updated:
			checkPixelColor(t, tr.lastTexture, 10, 10, color.White, "Pixel color mismatch after WhiteOp")
			checkPixelColor(t, tr.lastTexture, 400, 400, color.White, "Center pixel color mismatch after WhiteOp")
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for texture update after WhiteOp")
		}
	})

	t.Run("GreenBackground", func(t *testing.T) {
		l.Post(ResetOp{})
		l.Post(GreenOp{})
		l.Post(UpdateOp)

		select {
		case <-tr.updated:
			greenColor := color.RGBA{G: 255, A: 255}
			checkPixelColor(t, tr.lastTexture, 10, 10, greenColor, "Pixel color mismatch after GreenOp")
			checkPixelColor(t, tr.lastTexture, 400, 400, greenColor, "Center pixel color mismatch after GreenOp")
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for texture update after GreenOp")
		}
	})

	t.Run("WhiteBackgroundWithBlackRect", func(t *testing.T) {
		l.Post(ResetOp{})
		l.Post(WhiteOp{})
		l.Post(BgRectOp{X1: 0.25, Y1: 0.25, X2: 0.75, Y2: 0.75})
		l.Post(UpdateOp)

		select {
		case <-tr.updated:
			checkPixelColor(t, tr.lastTexture, 300, 300, color.Black, "Pixel inside black rect is not black")
			checkPixelColor(t, tr.lastTexture, 100, 100, color.White, "Pixel outside black rect is not white")
			checkPixelColor(t, tr.lastTexture, 700, 700, color.White, "Pixel outside black rect is not white")
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for texture update after BlackRectOp")
		}
	})

	t.Run("WhiteBackgroundWithFigure", func(t *testing.T) {
		l.Post(ResetOp{})
		l.Post(WhiteOp{})
		l.Post(FigureOp{X: 0.5, Y: 0.5})
		l.Post(UpdateOp)

		select {
		case <-tr.updated:
			figureColor := color.RGBA{R: 255, G: 255, B: 0, A: 255}

			checkPixelColor(t, tr.lastTexture, 400, 400, figureColor, "Center pixel is not figure color")
			checkPixelColor(t, tr.lastTexture, 350, 400, figureColor, "Pixel in horizontal part is not figure color")
			checkPixelColor(t, tr.lastTexture, 450, 400, figureColor, "Pixel in horizontal part is not figure color")
			checkPixelColor(t, tr.lastTexture, 400, 370, figureColor, "Pixel in vertical part is not figure color")
			checkPixelColor(t, tr.lastTexture, 400, 430, figureColor, "Pixel in vertical part is not figure color")
			checkPixelColor(t, tr.lastTexture, 100, 100, color.White, "Pixel outside figure is not white")
			checkPixelColor(t, tr.lastTexture, 700, 700, color.White, "Pixel outside figure is not white")
			checkPixelColor(t, tr.lastTexture, 500, 500, color.White, "Pixel at (500, 500) should be white after reset and setting white background")

		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for texture update after FigureOp")
		}
	})

	t.Run("MoveFigure", func(t *testing.T) {
		l.Post(ResetOp{})
		l.Post(WhiteOp{})
		l.Post(FigureOp{X: 0.5, Y: 0.5})
		l.Post(UpdateOp)

		select {
		case <-tr.updated:
			t.Log("Figure drawn in initial position for MoveFigure test")
		case <-time.After(time.Second):
			t.Fatal("Timeout drawing initial figure for MoveFigure test")
		}

		l.Post(MoveOp{X: 0.1, Y: 0.1})
		l.Post(UpdateOp)

		select {
		case <-tr.updated:
			figureColor := color.RGBA{R: 255, G: 255, B: 0, A: 255}

			checkPixelColor(t, tr.lastTexture, 480, 480, figureColor, "Pixel at new center is not figure color after move")
			checkPixelColor(t, tr.lastTexture, 400, 400, color.White, "Pixel at old center is not background after move")

		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for texture update after MoveOp")
		}
	})

	t.Run("ResetOperation", func(t *testing.T) {
		l.Post(ResetOp{})
		l.Post(GreenOp{})
		l.Post(BgRectOp{X1: 0.1, Y1: 0.1, X2: 0.9, Y2: 0.9})
		l.Post(FigureOp{X: 0.3, Y: 0.3})
		l.Post(UpdateOp)

		select {
		case <-tr.updated:
			t.Log("State populated before reset")
			greenColor := color.RGBA{G: 255, A: 255}
			checkPixelColor(t, tr.lastTexture, 50, 50, greenColor, "Pixel is not green before ResetOp")
			checkPixelColor(t, tr.lastTexture, 400, 400, color.Black, "Pixel is not black inside rect before ResetOp")

		case <-time.After(time.Second):
			t.Fatal("Timeout populating state before reset")
		}

		l.Post(ResetOp{})
		l.Post(UpdateOp)

		select {
		case <-tr.updated:
			checkPixelColor(t, tr.lastTexture, 100, 100, color.Black, "Pixel is not black after ResetOp")
			checkPixelColor(t, tr.lastTexture, 500, 500, color.Black, "Pixel is not black after ResetOp")
			checkPixelColor(t, tr.lastTexture, 400, 400, color.Black, "Center pixel is not black after ResetOp")

		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for texture update after ResetOp")
		}
	})
}