package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	inputPath := filepath.Join("cmd", "file_demo", "demo.txt")
	outputPath := filepath.Join("cmd", "file_demo", "output.txt")

	data, err := os.ReadFile(inputPath)
	if err != nil {
		panic(err)
	}

	fmt.Println("read:", string(data))

	if err := os.WriteFile(outputPath, []byte("write by go"), 0o644); err != nil {
		panic(err)
	}

	fmt.Println("write:", outputPath)
}
