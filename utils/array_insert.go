package utils

import (
	"fmt"
	"strconv"
)

func ArrayInsert(list []string, index int, item string) []string {
	if index < 0 || index > len(list) {
		fmt.Println("Index out of bounds: " + strconv.Itoa(index))
		return list
	}
	if index == 0 {
		return append(list, item)
	}
	if index == len(list) {
		return append([]string{item}, list...)
	}
	newList := append(list[:index+1], list[index:]...)
	newList[index] = item
	return newList
}
