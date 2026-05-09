package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Entry struct {
	Path string
	Size int64
}

func main() {
	topN := flag.Int("top", 0, "top N entries")
	minMB := flag.Int("min", 0, "minimum size in MB")
	flag.Parse()

	minSize := int64(*minMB) * 1024 * 1024

	target := "."
	if flag.NArg() > 0 {
		target = flag.Arg(0)
	}

	files, err := os.ReadDir(target)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	var entries []Entry

	ch := make(chan Entry, len(files))
	sem := make(chan struct{}, 8)
	for _, f := range files {
		path := filepath.Join(target, f.Name())

		sem <- struct{}{}

		go func(p string, f os.DirEntry) {
			defer func() { <-sem }()

			var size int64
			if f.IsDir() {
				size = dirSize(path)
			} else {
				info, err := f.Info()
				if err != nil {
					ch <- Entry{Path: p, Size: 0}
					return
				}
				size = info.Size()
			}

			ch <- Entry{
				Path: p,
				Size: size,
			}
		}(path, f)
	}

	for range files {
		e := <-ch
		entries = append(entries, e)
	}

	entries = filterByMinMB(entries, minSize)

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Size > entries[j].Size
	})

	if *topN > 0 && *topN < len(entries) {
		entries = entries[:*topN]
	}

	var total int64
	for _, f := range entries {
		total += f.Size
		fmt.Println(humanize(f.Size), shorten(f.Path))
	}
	fmt.Println("========================")
	fmt.Println("Total:", humanize(total))
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

func filterByMinMB(entries []Entry, minSize int64) []Entry {
	var filtered []Entry
	for _, e := range entries {
		if minSize > 0 && e.Size < minSize {
			continue
		}
		filtered = append(filtered, e)
	}
	return filtered
}

func shorten(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if strings.HasPrefix(path, home) {
		return "~" + strings.TrimPrefix(path, home)
	}

	return path
}
