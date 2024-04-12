package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	dir := flag.String("dir", "", "Count words from all files in the given directory")
	flag.Parse()
	// Calling the count function to count the number of words
	// received from the Standard Input and printing it out
	//fmt.Println(count(os.Stdin))
	var files []string

	if *dir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := filepath.WalkDir(*dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	}); err != nil {
		fmt.Fprintf(os.Stderr, "failed filepath.WalkDir: %s", err.Error())
	}

	for _, file := range files {
		data, err := os.Open(file)
		fmt.Println(data.Name())
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed Open: %s", err.Error())
		}
		//fmt.Println(file.Name())
		fmt.Println(count(data))
	}
}

func count(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	wc := 0
	for scanner.Scan() {
		wc++
	}
	return wc
}
