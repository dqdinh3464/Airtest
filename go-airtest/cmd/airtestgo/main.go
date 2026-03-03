package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"

	"github.com/AirtestProject/go-airtest/internal/android"
	"github.com/AirtestProject/go-airtest/internal/runner"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "screenshot":
		runScreenshot(os.Args[2:])
	case "find":
		runFind(os.Args[2:])
	case "touch":
		runTouch(os.Args[2:])
	case "tap":
		runTap(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}

func commonFlags(fs *flag.FlagSet) (adbPath, serial *string) {
	adbPath = fs.String("adb", "adb", "Path to adb executable")
	serial = fs.String("serial", "", "ADB device serial")
	return
}

func runScreenshot(args []string) {
	fs := flag.NewFlagSet("screenshot", flag.ExitOnError)
	adbPath, serial := commonFlags(fs)
	out := fs.String("out", "screen.png", "Output png path")
	_ = fs.Parse(args)

	client := android.NewClient(*adbPath, *serial)
	img, err := client.CaptureScreen()
	dieIfErr(err)

	f, err := os.Create(*out)
	dieIfErr(err)
	defer f.Close()
	dieIfErr(png.Encode(f, img))
	fmt.Printf("saved screenshot to %s\n", *out)
}

func runFind(args []string) {
	fs := flag.NewFlagSet("find", flag.ExitOnError)
	adbPath, serial := commonFlags(fs)
	tpl := fs.String("template", "", "Template image path")
	_ = fs.Parse(args)
	if *tpl == "" {
		die("--template is required")
	}

	eng := runner.NewEngine(android.NewClient(*adbPath, *serial))
	res, _, err := eng.FindTemplate(*tpl)
	dieIfErr(err)
	fmt.Printf("found at x=%d y=%d confidence=%.4f\n", res.X, res.Y, res.Confidence)
}

func runTouch(args []string) {
	fs := flag.NewFlagSet("touch", flag.ExitOnError)
	adbPath, serial := commonFlags(fs)
	tpl := fs.String("template", "", "Template image path")
	threshold := fs.Float64("threshold", 0.88, "Required confidence in [0..1]")
	_ = fs.Parse(args)
	if *tpl == "" {
		die("--template is required")
	}

	eng := runner.NewEngine(android.NewClient(*adbPath, *serial))
	res, err := eng.TouchTemplate(*tpl, *threshold)
	dieIfErr(err)
	fmt.Printf("touched template at x=%d y=%d confidence=%.4f\n", res.X, res.Y, res.Confidence)
}

func runTap(args []string) {
	fs := flag.NewFlagSet("tap", flag.ExitOnError)
	adbPath, serial := commonFlags(fs)
	x := fs.Int("x", 0, "X coordinate")
	y := fs.Int("y", 0, "Y coordinate")
	_ = fs.Parse(args)
	client := android.NewClient(*adbPath, *serial)
	dieIfErr(client.Tap(*x, *y))
	fmt.Printf("tapped at x=%d y=%d\n", *x, *y)
}

func usage() {
	fmt.Print(`airtestgo - Airtest-like Android UI automation in Go

Usage:
  airtestgo <command> [flags]

Commands:
  screenshot   Capture Android screenshot using adb
  find         Locate template in current screen
  touch        Find template then tap its center
  tap          Tap absolute coordinate
`)
}

func dieIfErr(err error) {
	if err != nil {
		die(err.Error())
	}
}

func die(msg string) {
	fmt.Fprintln(os.Stderr, "error:", msg)
	os.Exit(1)
}
