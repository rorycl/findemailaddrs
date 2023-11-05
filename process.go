package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

// addressMap keeps a map of unique addresses by lowercase email address
// addresses with isDoNotReply true are omitted
type addressMap map[string]address

// count returns the number of unique addresses
func (am addressMap) count() int {
	return len(am)
}

// dump writes the address map to an export tsf file
func (am addressMap) dump(f *os.File) error {
	// sort
	keys := []string{}
	for k := range am {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	writer := csv.NewWriter(f)
	writer.Comma = '\t'

	// write out addresses
	err := writer.Write([]string{"name", "email", "updated"})
	if err != nil {
		return err
	}
	for _, k := range keys {
		v := am[k]
		if err := writer.Write(v.stringSlice()); err != nil {
			return err
		}
	}
	return nil
}

var counter int = 0

// processFiles processes all the files found in walkerFindEmails
func processFiles(files []emailFile) (addressMap, error) {

	am := addressMap{}

	for _, f := range files {

		// make an enail record
		e := email{
			name: f.name,
			path: f.path,
		}

		// parse the email
		err := e.Parse()
		if err != nil {
			if errors.Is(err, parseIgnoreError) {
				fmt.Printf("skipping email parsing error from %s\n", e.path)
				if *verbose {
					fmt.Println("  ", err)
				}
				continue
			} else {
				return am, err
			}
		}

		// process message and add addresses in addressMap
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
	}
	return am, nil
}
