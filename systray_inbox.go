package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/getlantern/systray"
	"github.com/psanford/systray-inbox/icons"
	"gopkg.in/fsnotify.v1"
)

var directory string

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <directory_to_watch>\n", os.Args[0])
		os.Exit(1)
	}

	directory = os.Args[1]

	systray.Run(onReady, onExit)
}

func onExit() {
}

func onReady() {
	count := getCount(directory)
	fmt.Println("count", count)
	if count > 0 {
		systray.SetIcon(icons.BlueCircle)
	} else {
		systray.SetIcon(icons.WhiteCircle)
	}

	systray.SetTitle("Inbox")
	systray.SetTooltip("Inbox Count")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}

				count := getCount(directory)
				fmt.Println("count", count)
				if count > 0 {
					systray.SetIcon(icons.BlueCircle)
				} else {
					systray.SetIcon(icons.WhiteCircle)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Fatalf("Watcher error: %s", err)
			}
		}
	}()

	err = watcher.Add(directory)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func getCount(dir string) int {
	var count int
	dir = filepath.Clean(dir)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Walk error: %s\n", err)
			return err
		}

		if path != dir {
			count++
		}
		return nil
	})
	return count
}
