package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const (
	jpg  = ".jpg"
	jpeg = ".jpeg"
	png  = ".png"
)

var (
	directory string
	interval  time.Duration
)

func init() {
	log.SetFlags(0)
	home, err := os.UserHomeDir()
	if err == nil {
		home = filepath.Join(home, "Pictures")
	}
	flag.StringVar(&directory, "dir", home, "The directory containing the wallpapers.")
	flag.DurationVar(&interval, "interval", 30*time.Minute, "The interval for changing wallpaper. E.g. 300s, 5m, 1h. Minimum of 5m")
	flag.Parse()

	if err != nil || interval < 5*time.Minute {
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
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
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		dir, err := os.ReadDir(directory)
		if err != nil {
			log.Fatalln(err.Error())
		}

		imageFiles := make([]string, 0)
		for _, file := range dir {
			if isImageFile(filepath.Ext(file.Name())) {
				imageFiles = append(imageFiles, file.Name())
			}
		}
		if len(imageFiles) < 2 {
			log.Fatalf("%s needs more than 1 image for a slideshow.", directory)
		}

		imageIndicies := rand.Perm(len(imageFiles))
		for _, index := range imageIndicies {
			imagePath := filepath.Join(directory, imageFiles[index])
			if _, err := os.Stat(imagePath); err == nil {
				if err := setImage(imagePath); err != nil {
					log.Fatalln("Encountered an error when setting the background:", err)
				}
			} else {
				log.Printf("%s doesn't appear to exist. Re-scanning the directory for images.", imagePath)
				break
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
