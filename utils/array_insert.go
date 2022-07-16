package utils

import (
	"fmt"
	"strconv"
)

func ArrayInsert[V any](list []V, index int, item V) []V {
	if index < 0 || index > len(list) {
		fmt.Println("Index out of bounds: " + strconv.Itoa(index))
		return list
	}
	if index == 0 {
		return append(list, item)
	}
	if index == len(list) {
		return append([]V{item}, list...)
	}
	newList := append(list[:index+1], list[index:]...)
	newList[index] = item
	return newList
}
