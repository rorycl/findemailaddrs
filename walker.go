package main

import (
	"fmt"
	"os"
	"strings"
)

// emailFile points to a file containing emails on disk
type emailFile struct {
	name, path string
}

// files stores the found email files
var files []emailFile

// started checks the first input path is valid
var started bool

// check directory exists
func checkDirExists(dir string) bool {
	d, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer d.Close()

	info, err := d.Stat()
	if err != nil {
		return false
	}
	if !info.IsDir() {
		return false
	}
	return true
}

// walkerFindEmails is a custom file walker for eml files
func walkerFindEmails(path string, info os.FileInfo, err error) error {
	if !started {
		if !checkDirExists(path) {
			return fmt.Errorf("directory %s does not exist", path)
		}
		started = true
	}

	if *verbose && info.IsDir() {
		fmt.Println("processing directory:", info.Name())
	}

	if !info.IsDir() && strings.Contains(strings.ToLower(info.Name()), searchSuffix) {
		if *verbose {
			fmt.Println("   file:", info.Name())
		}
		files = append(files, emailFile{
			name: info.Name(),
			path: path,
		})
	}
	return nil
}
