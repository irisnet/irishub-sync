package helper

import (
	"fmt"
	"testing"
)

func TestArrayIterator_HasNext(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iterator := GetIterator(data)

	for iterator.HasNext() {
		data := iterator.Next().(int)
		if data == 1 || data == 3 {
			iterator.Remove()
		}
	}
	fmt.Println(iterator.Get())

	dataS := []string{"1", "2", "3", "4", "5"}
	iteratorS := GetIterator(dataS)

	for iteratorS.HasNext() {
		data := iteratorS.Next().(string)
		if data == "4" {
			iteratorS.Remove()
		}
	}
	fmt.Println(iteratorS.Get())

}
