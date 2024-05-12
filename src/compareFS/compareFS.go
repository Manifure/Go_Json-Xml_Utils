package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	oldFile := flag.String("old", "", "Path to the old snapshot file")
	newFile := flag.String("new", "", "Path to the new snapshot file")
	flag.Parse()

	if *oldFile == "" || *newFile == "" {
		fmt.Println("Please provide both --old and --new flags")
		os.Exit(1)
	}

	oldSnapshot := make(map[string]bool)

	// Читаем старый снапшот и сохраняем пути в мапу
	file, err := os.Open(*oldFile)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", *oldFile, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		path := scanner.Text()
		if path != "" {
			oldSnapshot[path] = true
		}
	}

	// Читаем новый снапшот и сравниваем с мапой старого снапшота
	file, err = os.Open(*newFile)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", *newFile, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		path := scanner.Text()
		if path != "" {
			if !oldSnapshot[path] {
				fmt.Printf("ADDED %s\n", path)
			}
			delete(oldSnapshot, path)
		}
	}

	// Выводим удаленные пути
	for path := range oldSnapshot {
		fmt.Printf("REMOVED %s\n", path)
	}
}
