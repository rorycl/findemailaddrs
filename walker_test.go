package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWalker(t *testing.T) {

	files := []string{}
	path := "testfiles/topdir/"

	// needed for walkerFindEmails
	fileChan = make(chan string)

	fChan := walker(path)
	for f := range fChan {
		files = append(files, f)
	}

	if got, want := len(files), 3; got != want {
		t.Errorf("got %d want %d files", got, want)
	}

	want := []string{
		"testfiles/topdir/dir 2/mail3.eml",
		"testfiles/topdir/dir1/email1.eml",
		"testfiles/topdir/dir1/email2.eml",
	}

	if diff := cmp.Diff(want, files); diff != "" {
		t.Errorf("walk mismatch (-want +got):\n%s", diff)
	}

	t.Logf("%#v\n", files)

}

// This test shows the caller needs to ascertain if the top level
// directory exists
func TestWalkerFail(t *testing.T) {

	Exiter = func(code int) {
		t.Fatalf("os.Exit was called with %d", code)
	}

	files := []string{}
	path := "testfiles/fail/"

	fChan := walker(path)
	for f := range fChan {
		files = append(files, f)
	}

	if len(files) != 0 {
		t.Error("expected 0 files")
	}
	t.Log("completed TestWalkerFail ok")

}
