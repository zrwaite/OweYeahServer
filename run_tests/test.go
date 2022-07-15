package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zrwaite/OweMate/terminal"
	"github.com/zrwaite/OweMate/tests"
	"github.com/zrwaite/OweMate/utils"
)

func main() {
	filter := false
	if len(os.Args) > 1 {
		filter = true
	}
	filterList := os.Args[1:]
	tests := []tests.Test{tests.BinarySearchTests, tests.LinearSearchTests}
	for _, test := range tests {
		if filter {
			found, _ := utils.LinearSearch(filterList, test.Name)
			if !found {
				continue
			}
		}
		success, err := test.Function()
		if !(success) {
			fmt.Println(terminal.Red + "FAIL: " + test.Name)
			log.Fatal(terminal.Reset + err.Error())
		} else {
			fmt.Println(terminal.Green + "PASS: " + test.Name)
			fmt.Print(terminal.Reset + "")
		}
	}
}
