package main

import (
	"testing"
)

func CheckQueryResult(t *testing.T) {

	// fmt.Println(GetContactData())
	want := 1
	if want == 0 {
		t.Errorf("want %q", want)
	}
}
