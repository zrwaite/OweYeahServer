package tests

import (
	"strconv"
	"testing"

	"github.com/zrwaite/OweYeah/utils"
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
	{[]string{"A", "B", "C", "D", "E", "F", "G"}, "G", true, 7},
}

func TestBinarySearch(t *testing.T) {
	for _, test := range tests {
		found, index := utils.StandardBinarySearch(test.List, test.Item)
		if found != test.Found || index != test.Index {
			PrintErr(t, "output did not match input. Expected: "+strconv.FormatBool(test.Found)+" at "+strconv.Itoa(test.Index)+" but got "+strconv.FormatBool(found)+" at "+strconv.Itoa(index))
		}
	}
}

func TestLinearSearch(t *testing.T) {
	for _, test := range tests {
		found, index := utils.StandardBinarySearch(test.List, test.Item)
		if found != test.Found || index != test.Index {
			PrintErr(t, "output did not match input. Expected: "+strconv.FormatBool(test.Found)+" at "+strconv.Itoa(test.Index)+" but got "+strconv.FormatBool(found)+" at "+strconv.Itoa(index))
		}
	}
}
