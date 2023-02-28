//go:build darwin

package main

import (
	"os/exec"
	"strconv"
)

func SetImage(imagePath string) error {
	return exec.Command("osascript", "-e", `tell application "System Events" to tell every desktop to set picture to `+strconv.Quote(imagePath)).Run()
}
