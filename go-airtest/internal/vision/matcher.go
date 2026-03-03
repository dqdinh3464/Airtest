package vision

import (
	"fmt"
	"image"
	"image/color"
)

// Result stores template matching output.
type Result struct {
	X          int
	Y          int
	Confidence float64
}

// MatchTemplate searches template in source image using normalized SAD over grayscale.
// Confidence is in [0..1], where 1 means perfect match.
func MatchTemplate(src image.Image, tpl image.Image) (Result, error) {
	sb := src.Bounds()
	tb := tpl.Bounds()

	tw := tb.Dx()
	th := tb.Dy()
	if tw <= 0 || th <= 0 {
		return Result{}, fmt.Errorf("template has invalid size")
	}
	if sb.Dx() < tw || sb.Dy() < th {
		return Result{}, fmt.Errorf("template size %dx%d is larger than source %dx%d", tw, th, sb.Dx(), sb.Dy())
	}

	tplGray := grayMatrix(tpl)
	best := Result{Confidence: -1}
	maxDiff := float64(255 * tw * th)

	for y := 0; y <= sb.Dy()-th; y++ {
		for x := 0; x <= sb.Dx()-tw; x++ {
			diff := 0.0
			for ty := 0; ty < th; ty++ {
				for tx := 0; tx < tw; tx++ {
					s := grayValue(src.At(sb.Min.X+x+tx, sb.Min.Y+y+ty))
					t := tplGray[ty][tx]
					if s > t {
						diff += float64(s - t)
					} else {
						diff += float64(t - s)
					}
				}
			}
			confidence := 1.0 - (diff / maxDiff)
			if confidence > best.Confidence {
				best = Result{X: x, Y: y, Confidence: confidence}
			}
		}
	}

	return best, nil
}

func grayMatrix(img image.Image) [][]uint8 {
	b := img.Bounds()
	out := make([][]uint8, b.Dy())
	for y := 0; y < b.Dy(); y++ {
		row := make([]uint8, b.Dx())
		for x := 0; x < b.Dx(); x++ {
			row[x] = grayValue(img.At(b.Min.X+x, b.Min.Y+y))
		}
		out[y] = row
	}
	return out
}

func grayValue(c color.Color) uint8 {
	r, g, b, _ := c.RGBA()
	// convert 16-bit range to 8-bit and apply luminance weights.
	r8 := float64(r >> 8)
	g8 := float64(g >> 8)
	b8 := float64(b >> 8)
	return uint8(0.299*r8 + 0.587*g8 + 0.114*b8)
}
