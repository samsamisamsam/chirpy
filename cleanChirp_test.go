package main

import "testing"

func TestCleanChirp(t *testing.T) {
	cases := []struct {
		Input    string
		Expected string
	}{
		{Input: "This is a kerfuffle opinion I need to share with the world", Expected: "This is a **** opinion I need to share with the world"},
		{Input: "A Sharbert is none to be messed with", Expected: "A **** is none to be messed with"},
		{Input: "Fornax my axx", Expected: "**** my axx"},
		{Input: "It is easy to spot a lorax fornaX", Expected: "It is easy to spot a lorax ****"},
		{Input: "It should be easy for a turtle to kerfuffle!", Expected: "It should be easy for a turtle to kerfuffle!"},
	}

	for _, test := range cases {
		output := cleanChirp(test.Input)
		if test.Expected != output {
			t.Errorf("Failed:\n Expected: %s\n Actual:  %s\n", test.Expected, output)
		}
	}
}
