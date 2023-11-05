package main

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWalker(t *testing.T) {

	started = false

	path := "testfiles/topdir/"

	err := filepath.Walk(path, walkerFindEmails)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(files), 3; got != want {
		t.Errorf("want %d files got %d files", got, want)
	}

	want := []emailFile{
		emailFile{name: "mail3.eml", path: "testfiles/topdir/dir 2/mail3.eml"},
		emailFile{name: "email1.eml", path: "testfiles/topdir/dir1/email1.eml"},
		emailFile{name: "email2.eml", path: "testfiles/topdir/dir1/email2.eml"},
	}

	if diff := cmp.Diff(want, files, cmp.AllowUnexported(emailFile{})); diff != "" {
		t.Errorf("walk mismatch (-want +got):\n%s", diff)
	}

	t.Logf("%#v\n", files)

}

func TestWalkerFail(t *testing.T) {

	started = false

	path := "testfiles/fail/"

	err := filepath.Walk(path, walkerFindEmails)
	if err == nil {
		t.Fatal("expected walk to fail")
	}

	t.Log(err)

}
