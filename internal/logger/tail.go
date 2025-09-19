package logger

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TailFile(filename string, lines int) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if lines > 0 {
		if err := printLastLines(file, lines); err != nil {
			return err
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(filename)
	if err != nil {
		return err
	}

	file.Seek(0, 2) // Seek to end
	scanner := bufio.NewScanner(file)

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				for scanner.Scan() {
					fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05"), scanner.Text())
				}
			}
		case err := <-watcher.Errors:
			return err
		}
	}
}

func printLastLines(file *os.File, n int) error {
	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > n {
			lines = lines[1:]
		}
	}
	
	for _, line := range lines {
		fmt.Println(line)
	}
	
	return scanner.Err()
}