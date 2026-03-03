package vision

import (
	"image"
	"image/color"
	"testing"
)

func TestMatchTemplate(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 8, 8))
	fill(src, color.RGBA{10, 10, 10, 255})

	// Put a bright 2x2 block at (3,4).
	for y := 4; y < 6; y++ {
		for x := 3; x < 5; x++ {
			src.Set(x, y, color.RGBA{250, 250, 250, 255})
		}
	}

	tpl := image.NewRGBA(image.Rect(0, 0, 2, 2))
	fill(tpl, color.RGBA{250, 250, 250, 255})

	res, err := MatchTemplate(src, tpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.X != 3 || res.Y != 4 {
		t.Fatalf("expected (3,4), got (%d,%d)", res.X, res.Y)
	}
	if res.Confidence < 0.99 {
		t.Fatalf("expected confidence > 0.99, got %.4f", res.Confidence)
	}
}

func fill(img *image.RGBA, c color.RGBA) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			img.Set(x, y, c)
		}
	}
}
