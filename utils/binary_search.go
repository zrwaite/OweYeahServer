package utils

import "github.com/zrwaite/OweMate/graph/model"

func GenericBinarySearch[V *model.User | GTLTType](list []V, item V, gt func(a, b V) bool, lt func(a, b V) bool) (found bool, index int) {
	startIndex := 0
	endIndex := len(list) - 1
	if len(list) == 0 {
		return false, 0
	}
	for {
		middleIndex := (startIndex + endIndex) / 2
		middleItem := list[middleIndex]
		// fmt.Println("Looking for " + item + ", found " + middleItem + " at " + strconv.Itoa(middleIndex) + " between " + strconv.Itoa(startIndex) + " and " + strconv.Itoa(endIndex))
		if lt(middleItem, item) {
			if startIndex == endIndex {
				return false, startIndex + 1
			}
			startIndex = middleIndex + 1
		} else if gt(middleItem, item) {
			if startIndex == endIndex {
				return false, startIndex
			}
			endIndex = middleIndex - 1
		} else {
			return true, middleIndex
		}
	}
}

type GTLTType interface {
	int | string
}

func CompareUserGT(a, b *model.User) bool {
	return a.Username > b.Username
}

func CompareUserLT(a, b *model.User) bool {
	return a.Username < b.Username
}

func UserBinarySearch(list []*model.User, user *model.User) (found bool, index int) {
	return GenericBinarySearch(list, user, CompareUserGT, CompareUserLT)
}

func StandardBinarySearch[V GTLTType](list []V, user V) (found bool, index int) {
	return GenericBinarySearch(list, user, func(a, b V) bool { return a > b }, func(a, b V) bool { return a < b })
}
