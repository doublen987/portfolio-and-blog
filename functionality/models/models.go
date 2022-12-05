package models

type ArrayItemInterface interface {
	Value() int
}

type ArrayItem int

func (arr ArrayItem) Value() int {
	return int(arr)
}
