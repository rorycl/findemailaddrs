package main

import (
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// integration test
func TestProcess(t *testing.T) {

	expected := []string{
		"example@test.net",
		"test3@test.net",
		"test4@test.net",
		"test5@test.net",
		"xxxx@test.net",
		"xxxxxxx@gmail.com",
		"xxxxyyy@gmail.com",
		"zparisian@testanother.com",
	}

	// needed for walkerFindEmails
	fileChan = make(chan string)

	// 1. get files to process
	fChan := walker("testfiles/topdir")

	// 2. process files
	emailChan, errorChan := processEmail(fChan)

	// 3. launch and complete digester
	am := addressMap{}
	done := make(chan struct{})
	go func() {
		var err error
		am, err = processUniqueEmails(emailChan, errorChan)
		if err != nil {
			fmt.Println("process error:", err)
			os.Exit(1)
		}
		done <- struct{}{}
	}()
	<-done

	keys := []string{}
	for k := range am {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	got, want := keys, expected
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("addresses mismatch (-want +got):\n%s", diff)
	}
}
