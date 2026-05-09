package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	target := "."

	files, err := os.ReadDir(target)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	for _, f := range files {
		path := filepath.Join(target, f.Name())

		if f.IsDir() {
			fmt.Println("[DIR]", path)
		} else {
			info, err := f.Info()
			if err != nil {
				continue
			}
			fmt.Println(info.Size(), path)
		}
	}
}
