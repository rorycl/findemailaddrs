// message
package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/mnako/letters"
)

// email represents the data file and some content relating to an email
// message file
type email struct {
	name  string // normalised name
	path  string // full path to file
	addrs []address
	date  time.Time
	err   error
}

// address is an email address with optional name
type address struct {
	name         string
	email        string
	seen         int
	date         time.Time
	isDoNotReply bool // is a "do not reply" address
}

// addressStringSlice is a slice of string for outputting to tab
// separated format
func (a *address) stringSlice() []string {
	return []string{a.name, a.email, a.date.Format("2006-01-02")}
}

// String is a string representation of an email message
func (e email) String() string {
	r := fmt.Sprintf("%s (%s) :", e.name, e.path)
	for _, a := range e.addrs {
		if a.name == "" {
			r += fmt.Sprintf(" <%s>", a.email)
		} else {
			r += fmt.Sprintf(" %s <%s>", a.name, a.email)
		}
	}
	return r
}

// processEmail processes email from paths provided by fileChan
func processEmail(fileChan <-chan string) (<-chan email, <-chan email) {

	emailChan := make(chan email)
	errorChan := make(chan email)

	go func() {
		defer close(emailChan)
		defer close(errorChan)
		var wg sync.WaitGroup
		for f := range fileChan {
			file := f
			wg.Add(1)
			go func(file string) {
				defer wg.Done()
				e := email{path: file, name: path.Base(file)}
				e.err = e.Parse()
				if e.err != nil {
					errorChan <- e
				} else {
					emailChan <- e
				}
			}(file)
		}
		wg.Wait()
	}()
	return emailChan, errorChan
}

// parseIgnoreError is an email parsing error that can be ignored
var parseIgnoreError error = errors.New("handled parsing error")

// parse an email message, catching errors
func (e *email) Parse() error {
	if e == nil {
		return errors.New("nil email provided")
	}
	if e.path == "" {
		return errors.New("empty path provided")
	}
	f, err := os.Open(e.path)
	if err != nil {
		return fmt.Errorf("open err %w", err)
	}

	m, err := letters.ParseEmail(f)
	if err != nil {
		if strings.Contains(err.Error(), "letters.parsers.parseAddressListHeader") {
			return errors.Join(parseIgnoreError, fmt.Errorf("parsing err %s %w", e.path, err))
		}
		return fmt.Errorf("parsing err %s %w", e.path, err)
	}

	e.date = m.Headers.Date
	allAddresses := m.Headers.From
	allAddresses = append(allAddresses, m.Headers.To...)
	allAddresses = append(allAddresses, m.Headers.Cc...)
	for _, a := range allAddresses {
		addr := address{
			name:  a.Name,
			email: a.Address,
			date:  e.date,
		}
		if strings.Contains(strings.ToLower(addr.email), "donotreply") {
			continue
		}
		if strings.Contains(strings.ToLower(addr.name), "undisclosed") {
			continue
		}
		e.addrs = append(e.addrs, addr)
	}
	return nil
}
