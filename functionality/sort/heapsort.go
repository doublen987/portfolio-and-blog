package sort

import "github.com/doublen987/Projects/MySite/server/functionality/models"

func Heapsort(array *[]models.ArrayItem) {
	n := len(*array)
	buildMaxHeap(array, n)
	for i := n; i > 0; i-- {
		tmp := (*array)[1]
		(*array)[1] = (*array)[i]
		(*array)[i] = tmp
		heapify(array, 1, n)
	}
}

func buildMaxHeap(array *[]models.ArrayItem, n int) {
	for i := n / 2; i > 1; i-- {
		heapify(array, i, n)
	}
}

func heapify(arrayref *[]models.ArrayItem, i int, n int) {
	left := 2 * i
	right := 2*i + 1

	array := (*arrayref)

	var max int

	if left <= n && array[left].Value() > array[i].Value() {
		max = left
	} else {
		max = i
	}

	if right <= n && array[right].Value() > array[max].Value() {
		max = right
	}

	if max != i {
		tmp := array[i]
		array[i] = array[max]
		array[max] = tmp
		heapify(&array, max, n)
	}
}
