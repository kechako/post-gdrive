package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kechako/post-gdrive/gdrive"
	isatty "github.com/mattn/go-isatty"

	"golang.org/x/net/context"
)

var isTerminal bool

func init() {
	isTerminal = isatty.IsTerminal(os.Stdout.Fd())
}

func _main() (int, error) {
	folderId := "root"
	filePath, err := os.Getwd()
	if err != nil {
		return 1, err
	}
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}
	if len(os.Args) > 2 {
		folderId = os.Args[2]
	}

	ctx := context.Background()

	g, err := gdrive.New(ctx)
	if err != nil {
		return 2, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return 3, err
	}
	defer file.Close()

	var size int64 = 0
	if info, err := file.Stat(); err == nil {
		size = info.Size()
	}

	f, err := g.UploadFile(filepath.Base(filePath), []string{folderId}, file, func(current, total int64) {
		printProgeress(current, size)
	})
	fmt.Println()
	if err != nil {
		return 5, err
	}

	fmt.Printf("Upload sucessful. [File Id: %s]\n", f.Id)

	return 0, nil
}

func printProgeress(current, size int64) {
	if isTerminal {
		if size == 0 {
			fmt.Printf("\r%d", current)
		} else {
			fmt.Printf("\r%d / %d (%d %%)", current, size, current*100/size)
		}
	} else {
		if size == 0 {
			fmt.Printf("%d\n", current)
		} else {
			fmt.Printf("%d / %d (%d %%)\n", current, size, current*100/size)
		}
	}
}

func main() {
	if exitStatus, err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "[Error] %v\n", err)
		os.Exit(exitStatus)
	}
}
