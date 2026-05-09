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
			size := dirSize(path)
			fmt.Println(size, path)
		} else {
			info, err := f.Info()
			if err != nil {
				continue
			}
			fmt.Println(info.Size(), path)
		}
	}
}

func dirSize(path string) int64 {
	var size int64

	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size
}
