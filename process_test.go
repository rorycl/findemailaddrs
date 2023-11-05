package main

import (
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

	am, err := processFiles(
		[]emailFile{
			emailFile{
				name: "mail1.eml",
				path: "testfiles/topdir/dir1/email1.eml",
			},
			emailFile{
				name: "mail2.eml",
				path: "testfiles/topdir/dir1/email2.eml",
			},
			emailFile{
				name: "mail3.eml",
				path: "testfiles/topdir/dir 2/mail3.eml",
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
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
