//go:build linux

package main

import (
	"errors"
	"os/exec"
	"strconv"
)

func SetImage(imagePath string) error {
	imageUri := strconv.Quote("file://" + imagePath)

	// Set both since we don't know if the desktop is a dark or light theme
	err1 := exec.Command("gsettings", args("picture-uri-dark", imageUri)...).Run()
	err2 := exec.Command("gsettings", args("picture-uri", imageUri)...).Run()
	return errors.Join(err1, err2)
}

func args(theme string, imageUri string) []string {
	return []string{"set", "org.gnome.desktop.background", theme, imageUri}
}
