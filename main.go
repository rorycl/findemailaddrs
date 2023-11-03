package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sg3des/eml"
)

// EML represents the data file and some content relating to an eml
// email file
type EML struct {
	name  string // normalised name
	path  string // full path to file
	addrs []address
	date  time.Time
}

// address is an email address with optional name
type address struct {
	name         string
	email        string
	isDoNotReply bool // is a "do not reply" address
}

// addressMap keeps a map of unique addresses by lowercase email address
// addresses with isDoNotReply true are omitted
type addressMap map[string]address

func (am addressMap) dump(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// sort
	keys := []string{}
	for k := range addressMap {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	// write out addresses
	_, err = f.Write("name,email\n")
	if err != nil {
		return err
	}
	for _, k := range keys {
		v := am[k]
		_, err := f.Write(fmt.Sprintf("%s,%s\n", v.name, v.email))
		if err != nil {
			return err
		}
	}
	return nil

}

// String is a string representation of an EML message
func (e EML) String() string {
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

var emlParseIgnoreError error = errors.New("non fatal eml parsing error")

// parse an EML message using the eml module, catching errors
func (e *EML) Parse() error {
	if e == nil {
		return errors.New("nil eml")
	}
	if e.path == "" {
		return errors.New("empty path provided")
	}
	f, err := os.Open(e.path)
	if err != nil {
		return fmt.Errorf("open err %w", err)
	}
	c, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("reading err %w", err)
	}
	m, err := eml.Parse(c)
	if err != nil {
		if strings.Contains(err.Error(), "multipart specified without boundary") {
			return emlParseIgnoreError
		}
		if strings.Contains(err.Error(), "invalid simpleAddr") {
			return emlParseIgnoreError
		}
		// this error condition is introduced at line 135 of
		// eml/address.go
		if strings.Contains(err.Error(), "invalid token length") {
			return emlParseIgnoreError
		}
		return fmt.Errorf("eml err %w", err)
	}
	e.date = m.Date
	allAddresses := m.From
	allAddresses = append(allAddresses, m.To...)
	allAddresses = append(allAddresses, m.Cc...)
	for _, a := range allAddresses {
		addr := address{
			name:  a.Name(),
			email: a.Email(),
		}
		if strings.Contains(strings.ToLower(addr.email), "donotreply") {
			e.addrs = append(e.addrs, addr)
		}
		e.addrs = append(e.addrs, addr)
	}
	return nil
}

func main() {

	directory := "/vols/jonathan/"

	// emls := []EML{}
	am := addressMap{}

	counter := 0
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {

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
			// fmt.Println(e)
			// emls = append(emls, e)
			// process
			counter++
			return nil
		}
		return nil
	})
	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Println("counter", counter)
}
