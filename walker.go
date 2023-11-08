package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// fileChan is a chan of file paths
var fileChan = make(chan string)

// Exiter indirect os.Exit
var Exiter = os.Exit

// walker generates a go routine
func walker(directory string) <-chan string {
	go func() {
		defer close(fileChan)
		err := filepath.Walk(directory, walkerFindEmails)
		if err != nil {
			fmt.Println("walk error:", err)
			Exiter(1)
		}
	}()
	return fileChan
}

// walkerFindEmails is a custom file walker for eml files
func walkerFindEmails(path string, info os.FileInfo, err error) error {

	if *verbose && info.IsDir() {
		fmt.Println("processing directory:", info.Name())
	}

	if !info.IsDir() && strings.Contains(strings.ToLower(info.Name()), searchSuffix) {
		if *verbose {
			fmt.Println("   file:", info.Name())
		}
		fileChan <- path // fileChan is declared at package level
	}
	return nil
}
