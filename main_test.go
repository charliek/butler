package main

import "testing"

func TestExample(t *testing.T) {
	var i = 2
	if i != 2 {
		t.Errorf("%v != %v", i, 2)
	}
}
