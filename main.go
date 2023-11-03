package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var counter int
var am = addressMap{}

func main() {

	directory := "/vols/jonathan/"

	counter := 0
	err := filepath.Walk(directory, walkerEML)
	if err != nil {
		fmt.Println("error: ", err)
	}
	err = am.dump("/tmp/outputX.csv")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("counter", counter)
}
