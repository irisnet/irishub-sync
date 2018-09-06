package helper

import (
	"fmt"
	"testing"
)

func TestArrayIterator_HasNext(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	iterator := NewIntIterator(&data)

	for iterator.HasNext() {
		data := iterator.Next().(int)
		if data == 5 {
			iterator.Remove()
		}
	}
	fmt.Println(iterator.Get())

}
