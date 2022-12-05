package sort

import "github.com/doublen987/Projects/MySite/server/functionality/models"

func Quicksort(array []models.ArrayItemInterface) []models.ArrayItemInterface {
	return quicksortRecursive(array, 0, len(array)-1)
}

func quicksortRecursive(array []models.ArrayItemInterface, low int, high int) []models.ArrayItemInterface {
	if low < high {
		var pivot int
		array, pivot = quicksortPartition(array, low, high)
		array = quicksortRecursive(array, low, pivot-1)
		array = quicksortRecursive(array, pivot+1, high)
	}
	return array
}

func quicksortPartition(array []models.ArrayItemInterface, low int, high int) ([]models.ArrayItemInterface, int) {
	pivot := array[high].Value()
	leftwall := low

	for i := low; i < high; i++ {
		if array[i].Value() < pivot {
			tmp := array[i]
			array[i] = array[leftwall]
			array[leftwall] = tmp
			leftwall = leftwall + 1
		}
	}

	tmp := array[high]
	array[high] = array[leftwall]
	array[leftwall] = tmp

	return array, leftwall
}
