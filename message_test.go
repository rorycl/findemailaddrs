package main

import (
	"fmt"
	"testing"
)

func TestMessages(t *testing.T) {

	testCases := []struct {
		name string
		path string
		want string
		err  bool
	}{
		{
			name: "mail1.eml",
			path: "testfiles/topdir/dir1/email1.eml",
			want: "mail1.eml (testfiles/topdir/dir1/email1.eml) : Xxxx Xxxxxxxxxxxxxx <xxxxxxx@gmail.com> Xxxx Xxxxxxxxxxxxxx <xxxx@test.net> <example@test.net>",
			err:  false,
		},
		{
			name: "mail2.eml",
			path: "testfiles/topdir/dir1/email2.eml",
			want: "mail2.eml (testfiles/topdir/dir1/email2.eml) : xxxx Xxxxxxxxxxxxxx <xxxxxxx@gmail.com> xxxx Xxxxxxxxxxxxxx <xxxx@test.net> <test3@test.net> Zanzibar Parisian <zparisian@testanother.com>",
			err:  false,
		},
		{
			name: "mail3.eml",
			path: "testfiles/topdir/dir 2/mail3.eml",
			want: "mail3.eml (testfiles/topdir/dir 2/mail3.eml) : Xxxx Xxxxxxxxxxxxxx <xxxxyyy@gmail.com> Xxxx Xxxxxxxxxxxxxx <xxxx@test.net> <test4@test.net> <test5@test.net>",
			err:  false,
		},
		{
			name: "fail.eml",
			path: "testfiles/topdir/dir3/fail.eml",
			want: "",
			err:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			e := email{name: tc.name, path: tc.path}
			err := e.Parse()
			if err != nil && tc.err == false {
				t.Fatal(err)
			}
			if err != nil && tc.err == true {
				return
			}
			if got, want := e.String(), tc.want; got != want {
				t.Errorf("got  %s\nwant %s", got, want)
				fmt.Println(got)
			}

		})
	}
}
