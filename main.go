package main

import (
	"fmt"
	"os"
)

func main() {
	currentPath := "."

	if len(os.Args) > 1 {
		currentPath = os.Args[1]
	}

	fmt.Println("scan target:", currentPath)
}
