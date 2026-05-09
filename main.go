package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/atolix/duct/internal/scan"
	"github.com/atolix/duct/internal/util"
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

	entries = filterByMinSize(entries, minSize)
	sortEntries(entries)
	entries = takeTopN(entries, *topN)

	printEntries(entries)
}

func sortEntries(entries []scan.Entry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Size > entries[j].Size
	})
}

func takeTopN(entries []scan.Entry, topN int) []scan.Entry {
	if topN > 0 && topN < len(entries) {
		entries = entries[:topN]
	}
	return entries
}

func filterByMinSize(entries []scan.Entry, minSize int64) []scan.Entry {
	var filtered []scan.Entry
	for _, e := range entries {
		if minSize > 0 && e.Size < minSize {
			continue
		}
		filtered = append(filtered, e)
	}
	return filtered
}

func printEntries(entries []scan.Entry) {
	var total int64
	for _, f := range entries {
		total += f.Size
		fmt.Println(util.Humanize(f.Size), util.Shorten(f.Path))
	}
	fmt.Println("========================")
	fmt.Println("Total:", util.Humanize(total))
}
