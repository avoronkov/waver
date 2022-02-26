package watch

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func OnFileUpdate(filename string, action func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("Cannot detect abs path of %v: %v", filename, err)
	}
	dir := filepath.Dir(absPath)
	log.Printf("Starting watching file changes: %v (%v)", absPath, dir)
	if err := watcher.Add(dir); err != nil {
		return err
	}
	go func() {
	L:
		for {
			select {
			case event, ok := <-watcher.Events:
				// log.Printf("File event: %v (%v)", event, event.Name)
				if !ok {
					break L
				}
				if event.Name != absPath {
					break
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					action()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					break L
				}
				log.Printf("File watcher error: %v", err)

			}
		}
		log.Printf("File watcher is stopped.")
	}()

	return nil
}
