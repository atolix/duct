package scan

import (
	"os"
	"path/filepath"
	"sync"
)

type Entry struct {
	Path string
	Size int64
}

func ScanDir(target string) ([]Entry, error) {
	files, err := os.ReadDir(target)
	if err != nil {
		return nil, err
	}

	var entries []Entry

	ch := make(chan Entry, len(files))
	sem := make(chan struct{}, 8)
	var wg sync.WaitGroup
	for _, f := range files {
		path := filepath.Join(target, f.Name())

		sem <- struct{}{}
		wg.Add(1)

		go func(p string, f os.DirEntry) {
			defer wg.Done()
			defer func() { <-sem }()

			var size int64
			if f.IsDir() {
				size = dirSize(p)
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

	go func() {
		wg.Wait()
		close(ch)
	}()

	for e := range ch {
		entries = append(entries, e)
	}

	return entries, nil
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

