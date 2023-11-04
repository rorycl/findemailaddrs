package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// map of addresses to address struct
var am = addressMap{}

// default search file extension
var searchSuffix = ".eml"

// flags
var (
	directory = flag.String("d", "", "path to directory to start eml file search")
	output    = flag.String("o", "", "file to save output")
	suffix    = flag.String("s", searchSuffix, "file suffix to search for")
	verbose   = flag.Bool("v", false, "verbose")
)

var usage = `%s -d <directory> -o <output> [-v] [-s ".ext"]

Look for files with the default "%s" suffix or that optionally provided
with the "-s" flag in the directory rooted at <directory> and extract
the email addresses and associated names (where available) to <output>
in tab separated format.

Provide the -v flag for verbose output.

Options:
`

func main() {

	flag.Usage = func() {
		fmt.Printf(usage, os.Args[0], searchSuffix)
		flag.PrintDefaults()
	}

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

	if *suffix != searchSuffix {
		searchSuffix = *suffix
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

	err = filepath.Walk(*directory, walkerFindEmails)
	if err != nil {
		fmt.Println("error: ", err)
	}
	err = am.dump(outputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// stats
	if *verbose {
		fmt.Println("counter", counter)
		fmt.Println("unique addresses", am.count())
	}
}
