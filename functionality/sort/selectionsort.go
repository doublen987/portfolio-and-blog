package sort

import "github.com/doublen987/Projects/MySite/server/functionality/models"

func SelectionSort(array []models.ArrayItemInterface) []models.ArrayItemInterface {
	for i := 0; i < len(array); i++ {
		minvalue := i
		for j := i; j < len(array); j++ {
			if array[i].Value() < array[j].Value() {
				minvalue = j
			}
		}
		tmp := array[i]
		array[i] = array[minvalue]
		array[minvalue] = tmp
	}
	return array
}
