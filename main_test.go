package main

import "testing"

func TestAdd(t *testing.T) {
	params := []struct {
		a, b, c int
	}{
		{1, 2, 3},
		{2, 3, 5},
		{3, 4, 7},
		{5, 6, 11},
		{-1, 1, 0},
	}

	for _, param := range params {
		actural := add(param.a, param.b)
		if actural != param.c {
			t.Errorf("add(%d,%d) = %d; expected %d", param.a, param.b, actural, param.c)
		}
	}

}
