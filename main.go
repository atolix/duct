package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/atolix/duct/internal/scan"
)

func main() {
	topN := flag.Int("top", 0, "top N entries")
	minMB := flag.Int("min", 0, "minimum size in MB")
	flag.Parse()

	minSize := int64(*minMB) * 1024 * 1024

	target := "."
	if flag.NArg() > 0 {
		target = flag.Arg(0)
	}

	entries, err := scan.ScanDir(target)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	entries = filterByMinMB(entries, minSize)
	sortEntries(entries)

	if *topN > 0 && *topN < len(entries) {
		entries = entries[:*topN]
	}

	printEntries(entries)
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

func filterByMinMB(entries []scan.Entry, minSize int64) []scan.Entry {
	var filtered []scan.Entry
	for _, e := range entries {
		if minSize > 0 && e.Size < minSize {
			continue
		}
		filtered = append(filtered, e)
	}
	return filtered
}

func sortEntries(entries []scan.Entry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Size > entries[j].Size
	})
}

func shorten(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if rest, ok := strings.CutPrefix(path, home); ok {
		return "~" + rest
	}

	return path
}

func printEntries(entries []scan.Entry) {
	var total int64
	for _, f := range entries {
		total += f.Size
		fmt.Println(humanize(f.Size), shorten(f.Path))
	}
	fmt.Println("========================")
	fmt.Println("Total:", humanize(total))
}
