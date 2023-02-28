//go:build windows

package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-systemparametersinfoa
var (
	user32Dll           = windows.NewLazyDLL("user32.dll")
	procSystemParamInfo = user32Dll.NewProc("SystemParametersInfoW")
)

const (
	spiSetDeskWallpaper = 0x0014
	uiParam             = 0x0000
	spifUpdateIniFile   = 0x001A
)

func SetImage(imagePath string) error {
	imagePathPtr, _ := windows.UTF16PtrFromString(imagePath)
	procSystemParamInfo.Call(spiSetDeskWallpaper, uiParam, uintptr(unsafe.Pointer(imagePathPtr)), spifUpdateIniFile)
	return nil // TODO error handling
}
