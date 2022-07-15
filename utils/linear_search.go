package utils

func LinearSearch(list []string, item string) (found bool, index int) {
	for index, listItem := range list {
		if listItem == item {
			return true, index
		} else if listItem > item {
			return false, index - 1
		}
	}
	return false, len(list)
}
