package sort

import "github.com/doublen987/Projects/MySite/server/functionality/models"

func BubbleSort(array []models.ArrayItemInterface) []models.ArrayItemInterface {
	n := len(array)
	for i := 1; i < n; i++ {
		for j := 0; j < n-1; j++ {
			if array[i].Value() > array[j+1].Value() {
				tmp := array[i]
				array[i] = array[j+1]
				array[j+1] = tmp
			}
		}
	}
	return array
}
