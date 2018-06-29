package trim

import (
	"fmt"
	"strings"
	"testing"
)

func Example() {
	rset, err := NewRuneSet("a-zA-Z")
	if err != nil {
		panic(err)
	}
	fmt.Println(rset.Trim("Hello, World!"))
	// Output:
	// HelloWorld
}

func TestTrim(t *testing.T) {
	tcs := []struct {
		rset string
		in   string
		out  string
	}{
		{
			rset: "a-z",
			in:   "Hello, World!",
			out:  "elloorld",
		},
	}
	for _, tc := range tcs {
		rset, err := NewRuneSet(tc.rset)
		if err != nil {
			t.Fatalf("couldn't compile %s", tc.rset)
		}
		out := rset.Trim(tc.in)
		if out != tc.out {
			t.Errorf("expected %s but got %s", tc.out, out)
		}
	}
}

func TestMakeRuneSet(t *testing.T) {
	tcs := []struct {
		rset string
		in   string
		out  string
		err  string
	}{
		{
			rset: "a-z",
			in:   "bfk",
			out:  "1A他",
		},
		{
			rset: "uaeio",
			in:   "aiou",
			out:  "1bczA",
		},
		{
			rset: "a-zA-Z0-9+/",
			in:   "kK4+/",
			out:  "_.-",
		},
		{
			rset: "abc-",
			in:   "b-",
			out:  "d_",
		},
		{
			rset: "a-z他-我",
			in:   "lq你",
			out:  "A_人",
		},
		{
			rset: "z-a",
			err:  "invalid range z-a",
		},
		{
			rset: "ABCz-a",
			err:  "invalid range z-a",
		},
		{
			rset: "a- ",
			err:  "includes surrogates",
		},
		{
			rset: "a-oi-t",
			err:  "range i-t overlaps with a-o",
		},
		{
			rset: "fa-o",
			err:  "range a-o overlaps with f-f",
		},
		{
			rset: "a-za-c",
			err:  "range a-c overlaps with a-z",
		},
		{
			rset: "a-ca-z",
			err:  "range a-z overlaps with a-c",
		},
	}
	for _, tc := range tcs {
		rset, err := NewRuneSet(tc.rset)
		if tc.err != "" {
			if err == nil || !strings.Contains(err.Error(), tc.err) {
				t.Errorf("expected error %v to contain: %s", err, tc.err)
			}
			continue
		}
		for _, r := range tc.in {
			if !rset.Includes(r) {
				t.Errorf("expected %s to include %s", tc.rset, string(r))
				t.Fatalf("expected %s to include %s: %#v", tc.rset, string(r), rset)
			}
		}
		for _, r := range tc.out {
			if rset.Includes(r) {
				t.Errorf("didn't expect %s to include %s", tc.rset, string(r))
			}
		}
	}
}
