package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// map of addresses to address struct
var am = addressMap{}

// flags
var directory = flag.String("d", "", "path to directory to start eml file search")
var output = flag.String("o", "", "file to save output")

func main() {

	flag.Parse()
	if directory == nil {
		fmt.Printf("no directory provided")
		flag.Usage()
		os.Exit(1)
	}
	if output == nil {
		fmt.Printf("no output file provided")
		flag.Usage()
		os.Exit(1)
	}

	// check directory exists
	d, err := os.Open(*directory)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	info, err := d.Stat()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Println("path to directory is not a directory")
	}

	// try and open output file
	outputFile, err := os.Create(*output)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer outputFile.Close()

	err = filepath.Walk(*directory, walkerEML)
	if err != nil {
		fmt.Println("error: ", err)
	}
	err = am.dump(outputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("counter", counter)
}
