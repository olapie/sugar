package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s {modelFilename}", os.Args[0])
		return
	}
	Generate(os.Args[1])
}
