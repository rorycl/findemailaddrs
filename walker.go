package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var counter int = 0

// walkerEML is a custom file walker for eml files
func walkerEML(path string, info os.FileInfo, err error) error {

	if !info.IsDir() && strings.Contains(strings.ToLower(info.Name()), ".eml") {
		e := EML{
			name: info.Name(),
			path: path,
		}
		err := e.Parse()
		if err != nil {
			if errors.Is(err, emlParseIgnoreError) {
				return nil
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// process eml and put addresses in addressMap
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
			} else {
				am[email] = a
			}
		}
		counter++
		return nil
	}
	return nil
}
