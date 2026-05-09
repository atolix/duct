package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type Entry struct {
	Path string
	Size int64
}

func main() {
	target := "."

	files, err := os.ReadDir(target)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	var entries []Entry

	for _, f := range files {
		path := filepath.Join(target, f.Name())

		if f.IsDir() {
			size := dirSize(path)
			entries = append(entries, Entry{
				Path: path,
				Size: size,
			})
		} else {
			info, err := f.Info()
			if err != nil {
				continue
			}
			size := info.Size()
			entries = append(entries, Entry{
				Path: path,
				Size: size,
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Size > entries[j].Size
	})

	for _, e := range entries {
		fmt.Println(humanize(e.Size), e.Path)
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

func humanize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(size)/float64(div), "KMGTPE"[exp])
}
