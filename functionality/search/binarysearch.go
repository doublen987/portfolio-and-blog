package search

import "github.com/doublen987/Projects/MySite/server/functionality/models"

func BinarySearch(array []models.ArrayItemInterface, target int) int {
	left := 0
	right := len(array) - 1

	for left <= right {
		mid := (right + left) / 2

		if array[mid].Value() == target {
			return mid
		} else if array[mid].Value() < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1
}
