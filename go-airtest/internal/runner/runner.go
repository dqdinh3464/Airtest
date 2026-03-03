package runner

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/AirtestProject/go-airtest/internal/android"
	"github.com/AirtestProject/go-airtest/internal/vision"
)

// Engine combines device and template matching features.
type Engine struct {
	Device *android.Client
}

func NewEngine(device *android.Client) *Engine {
	return &Engine{Device: device}
}

func LoadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open image %s: %w", path, err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode image %s: %w", path, err)
	}
	return img, nil
}

func (e *Engine) FindTemplate(templatePath string) (vision.Result, image.Image, error) {
	screen, err := e.Device.CaptureScreen()
	if err != nil {
		return vision.Result{}, nil, err
	}
	tpl, err := LoadImage(templatePath)
	if err != nil {
		return vision.Result{}, nil, err
	}
	res, err := vision.MatchTemplate(screen, tpl)
	if err != nil {
		return vision.Result{}, nil, err
	}
	return res, tpl, nil
}

func (e *Engine) TouchTemplate(templatePath string, threshold float64) (vision.Result, error) {
	res, tpl, err := e.FindTemplate(templatePath)
	if err != nil {
		return vision.Result{}, err
	}
	if res.Confidence < threshold {
		return res, fmt.Errorf("match confidence %.4f is lower than threshold %.4f", res.Confidence, threshold)
	}
	b := tpl.Bounds()
	cx := res.X + b.Dx()/2
	cy := res.Y + b.Dy()/2
	if err := e.Device.Tap(cx, cy); err != nil {
		return res, err
	}
	return res, nil
}
