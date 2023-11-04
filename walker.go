package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var counter int = 0

// walkerFindEmails is a custom file walker for eml files
func walkerFindEmails(path string, info os.FileInfo, err error) error {

	if *verbose && info.IsDir() {
		fmt.Println("processing directory:", info.Name())
	}

	if !info.IsDir() && strings.Contains(strings.ToLower(info.Name()), searchSuffix) {
		if *verbose {
			fmt.Println("   file:", info.Name())
		}
		e := email{
			name: info.Name(),
			path: path,
		}

		// the main parsing of the email occurs here
		err := e.Parse()
		if err != nil {
			if errors.Is(err, parseIgnoreError) {
				fmt.Printf("skipping email parsing error from %s\n", path)
				if *verbose {
					fmt.Println("  ", err)
				}
				return nil
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// process message and put addresses in addressMap
		for _, a := range e.addrs {
			if a.isDoNotReply {
				continue
			}
			email := strings.ToLower(a.email)
			existing, ok := am[email]
			if ok {
				existing.seen++
				if existing.name == "" && a.name != "" {
					existing.name = a.name
				}
				if existing.date.Before(a.date) {
					existing.date = a.date
				}
			} else {
				am[email] = a
			}
		}
		counter++
		return nil
	}
	return nil
}
