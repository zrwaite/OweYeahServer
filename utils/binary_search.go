package utils

func BinarySearch(list []string, item string) (found bool, index int) {
	startIndex := 0
	endIndex := len(list) - 1
	if len(list) == 0 {
		return false, 0
	}
	for {
		middleIndex := (startIndex + endIndex) / 2
		middleItem := list[middleIndex]
		// fmt.Println("Looking for " + item + ", found " + middleItem + " at " + strconv.Itoa(middleIndex) + " between " + strconv.Itoa(startIndex) + " and " + strconv.Itoa(endIndex))
		if middleItem < item {
			if startIndex == endIndex {
				return false, startIndex + 1
			}
			startIndex = middleIndex + 1
		} else if middleItem > item {
			if startIndex == endIndex {
				return false, startIndex
			}
			endIndex = middleIndex - 1
		} else {
			return true, middleIndex
		}
	}
}
