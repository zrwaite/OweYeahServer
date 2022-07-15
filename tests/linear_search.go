package tests

import (
	"errors"

	"github.com/zrwaite/OweMate/utils"
)

type SearchTest struct {
	List  []string
	Item  string
	Found bool
	Index int
}

var tests = []SearchTest{
	{[]string{"A", "B", "C", "D"}, "A", true, 0},
	{[]string{"A", "B", "C", "D", "E"}, "C", true, 2},
	{[]string{"A", "B", "C", "D", "F"}, "E", false, 4},
	{[]string{"B", "C", "D", "F"}, "A", false, 0},
	{[]string{"A", "B", "C", "D", "E", "F", "G"}, "H", false, 7},
	{[]string{"A", "B", "C", "D", "E", "F", "G"}, "G", true, 6},
}

var BinarySearchTests = Test{
	Name: "Binary Search Test",
	Function: func() (bool, error) {
		for _, test := range tests {
			found, index := utils.BinarySearch(test.List, test.Item)
			if found != test.Found || index != test.Index {
				return false, errors.New("output did not match input")
			}
		}
		return true, nil
	},
}

var LinearSearchTests = Test{
	Name: "Linear Search Test",
	Function: func() (bool, error) {
		for _, test := range tests {
			found, index := utils.BinarySearch(test.List, test.Item)
			if found != test.Found || index != test.Index {
				return false, errors.New("output did not match input")
			}
		}
		return true, nil
	},
}
