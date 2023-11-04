// message
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
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
}

// address is an email address with optional name
type address struct {
	name         string
	email        string
	seen         int
	date         time.Time
	colleague    bool
	isDoNotReply bool // is a "do not reply" address
}

// addressStringSlice is a slice of string for outputting to tab
// separated formate
func (a *address) stringSlice() []string {
	colBool := "false"
	if a.colleague {
		colBool = "true"
	}
	return []string{a.name, a.email, a.date.Format("2006-01-02"), colBool}
}

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
	err := writer.Write([]string{"name", "email", "updated", "colleague"})
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

var parseIgnoreError error = errors.New("handled parsing error")

// parse an email message using letter , catching errors
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
		if strings.Contains(strings.ToLower(addr.email), "ucl.ac.uk") {
			addr.colleague = true
		}
		e.addrs = append(e.addrs, addr)
	}
	return nil
}
