package android

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os/exec"
	"strconv"
)

// Client provides a minimal Android device control layer through adb.
type Client struct {
	AdbPath string
	Serial  string
}

func NewClient(adbPath, serial string) *Client {
	if adbPath == "" {
		adbPath = "adb"
	}
	return &Client{AdbPath: adbPath, Serial: serial}
}

func (c *Client) baseArgs() []string {
	args := make([]string, 0, 2)
	if c.Serial != "" {
		args = append(args, "-s", c.Serial)
	}
	return args
}

func (c *Client) run(args ...string) ([]byte, error) {
	full := append(c.baseArgs(), args...)
	cmd := exec.Command(c.AdbPath, full...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("adb %v failed: %w, output: %s", args, err, bytes.TrimSpace(out))
	}
	return out, nil
}

func (c *Client) Shell(args ...string) error {
	shellArgs := append([]string{"shell"}, args...)
	_, err := c.run(shellArgs...)
	return err
}

func (c *Client) Tap(x, y int) error {
	return c.Shell("input", "tap", strconv.Itoa(x), strconv.Itoa(y))
}

func (c *Client) Swipe(x1, y1, x2, y2, ms int) error {
	return c.Shell("input", "swipe", strconv.Itoa(x1), strconv.Itoa(y1), strconv.Itoa(x2), strconv.Itoa(y2), strconv.Itoa(ms))
}

func (c *Client) KeyEvent(keyCode int) error {
	return c.Shell("input", "keyevent", strconv.Itoa(keyCode))
}

// CaptureScreen captures Android screen using `adb exec-out screencap -p`.
func (c *Client) CaptureScreen() (image.Image, error) {
	out, err := c.run("exec-out", "screencap", "-p")
	if err != nil {
		return nil, err
	}
	// Some devices return CRLF line endings, normalize to valid PNG.
	out = bytes.ReplaceAll(out, []byte("\r\n"), []byte("\n"))
	img, err := png.Decode(bytes.NewReader(out))
	if err != nil {
		return nil, fmt.Errorf("decode screencap png: %w", err)
	}
	return img, nil
}
