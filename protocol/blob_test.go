package protocol

import (
	"testing"
	"testing/quick"
)

func TestBlobBothWays(t *testing.T) {
	f := func(s []string) []string {
		return s
	}
	g := func(s []string) []string {
		b := &Blob{blob: s}
		return DecodeBlob(b.Encode())
	}
	err := quick.CheckEqual(f, g, nil)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
