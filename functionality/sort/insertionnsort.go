package sort

import "github.com/doublen987/Projects/MySite/server/functionality/models"

func InsertionSort(array []models.ArrayItemInterface) []models.ArrayItemInterface {
	for i := 1; i < len(array); i++ {
		j := i
		for j > 0 && array[j-1].Value() > array[j].Value() {
			tmp := array[j]
			array[j] = array[j-1]
			array[j-1] = tmp
			j = j - 1
		}
	}
	return array
}
