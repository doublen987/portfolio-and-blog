package sort

import "github.com/doublen987/Projects/MySite/server/functionality/models"

func MergeSort(array []models.ArrayItemInterface) []models.ArrayItemInterface {
	n := len(array)

	if n == 1 {
		return array
	}
	arrayOne := array[:n/2]
	arrayTwo := array[n/2+1:]

	arrayOne = MergeSort(arrayOne)
	arrayTwo = MergeSort(arrayTwo)

	return merge(arrayOne, arrayTwo)

}

func merge(a []models.ArrayItemInterface, b []models.ArrayItemInterface) []models.ArrayItemInterface {
	c := []models.ArrayItemInterface{}

	for len(a) != 0 && len(b) != 0 {
		if a[0].Value() > b[0].Value() {
			c = append(c, b[0])
			b = b[1:]
		} else {
			c = append(c, a[0])
			a = a[1:]
		}
	}

	for len(a) != 0 {
		c = append(c, a[0])
		a = a[1:]
	}

	for len(b) != 0 {
		c = append(c, b[0])
		b = b[1:]
	}

	return c
}
