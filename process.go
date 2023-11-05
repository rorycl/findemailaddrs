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
	writer.Flush()
	return nil
}

var counter int = 0

// processUniqueEmails processes emails into a unique map
func processUniqueEmails(emailChan <-chan email, errorChan <-chan email) (addressMap, error) {

	am := addressMap{}

EMAIL:
	for {
		select {

		// read off errors from processEmail
		case e, ok := <-errorChan:
			if !ok {
				break EMAIL
			}
			if errors.Is(e.err, parseIgnoreError) {
				fmt.Printf("skipping email parsing error from %s\n", e.path)
				if *verbose {
					fmt.Println("  ", e.err)
				}
				continue
			} else {
				// unrecoverable error
				return am, e.err
			}

		// read off emails from processEmail
		case e, ok := <-emailChan:
			if !ok {
				break EMAIL
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
	}
	return am, nil
}
