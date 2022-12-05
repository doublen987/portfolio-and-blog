package sort

import (
	"testing"

	"github.com/doublen987/Projects/MySite/server/functionality/models"
)

func TestQuicksort(t *testing.T) {
	arr := []models.ArrayItem{3, 4, 12, 1, 3, 2, 45, 123, 14, 8}

	b := make([]models.ArrayItemInterface, len(arr))
	for i := range arr {
		b[i] = arr[i]
	}
	newArr := Quicksort(b)
	for _, item := range newArr {
		t.Log(item.Value())
	}
	lastElement := 0
	for _, element := range newArr {
		if element.Value() < lastElement {
			t.Fatal("Elements of the array are not in ascending order")
		}
		lastElement = element.Value()
	}
}
