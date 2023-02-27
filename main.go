package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

const usage = "Usage: bgslide /path/to/wallpapers intervalSeconds"

const (
	jpg  = ".jpg"
	jpeg = ".jpeg"
	png  = ".png"
)

func main() {
	log.SetFlags(0)
	ctx := cancelCtxOnSigterm(context.Background())
	run(ctx, SetImage)
}

func isImageFile(fileExtension string) bool {
	switch fileExtension {
	case jpg,
		jpeg,
		png:
		return true
	default:
		return false
	}
}

func cancelCtxOnSigterm(ctx context.Context) context.Context {
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-exitCh
		cancel()
	}()
	return ctx
}

func run(ctx context.Context, setImage func(imagePath string) error) {
	if len(os.Args) < 3 {
		log.Fatalln(usage)
	}

	wallpaperDir := os.Args[1]
	interval, err := strconv.ParseInt(os.Args[2], 10, 0)
	if err != nil {
		log.Fatalln(usage)
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	for {
		dir, err := os.ReadDir(wallpaperDir)
		if err != nil {
			log.Fatalln(err.Error())
		}

		imageFiles := make([]string, 0)
		for _, file := range dir {
			if isImageFile(filepath.Ext(file.Name())) {
				imageFiles = append(imageFiles, file.Name())
			}
		}

		imageIndicies := rand.Perm(len(imageFiles))
		for _, index := range imageIndicies {
			imagePath := filepath.Join(wallpaperDir, imageFiles[index])
			if err := setImage(imagePath); err != nil {
				log.Fatalln("Encountered an error when setting the background:", err)
			}

			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				return
			}
		}
	}
}
