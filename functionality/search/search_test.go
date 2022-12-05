package search

import (
	"testing"

	"github.com/doublen987/Projects/MySite/server/functionality/models"
)

func TestBinarySearch(t *testing.T) {
	arr := []models.ArrayItem{1, 2, 5, 12, 17, 25, 28, 33, 40, 49, 52, 56, 57, 60, 63, 67, 70}
	arrInterfaces := []models.ArrayItemInterface{}
	for _, element := range arr {
		arrInterfaces = append(arrInterfaces, element)
	}
	found := BinarySearch(arrInterfaces, 63)
	if found == -1 {
		t.Fatal("Could not find item")
	}
}
