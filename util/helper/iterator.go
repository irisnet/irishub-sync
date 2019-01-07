package helper

import (
	"reflect"
)

type Iterator interface {
	HasNext() bool
	Next() interface{}
	Remove()
	Get() interface{}
	Length() int
}

func GetIterator(data interface{}) Iterator {
	return &iterator{
		data: data,
	}
}

type iterator struct {
	current int
	data    interface{}
}

func (iterator *iterator) HasNext() bool {
	return iterator.current < iterator.Length()
}

func (iterator *iterator) Next() (v interface{}) {
	v = reflect.ValueOf(iterator.data).Index(iterator.current).Interface()
	iterator.current++
	return v
}

func (iterator *iterator) Remove() {
	var len = iterator.Length()
	frontData, afterData := reflect.ValueOf(iterator.data).Slice(0, iterator.current-1), reflect.ValueOf(iterator.data).Slice(iterator.current, len)
	newData := reflect.AppendSlice(frontData, afterData)
	iterator.data = newData.Interface()
	iterator.current -= 1
}

func (iterator *iterator) Get() interface{} {
	return iterator.data
}

func (iterator *iterator) Length() int {
	return reflect.ValueOf(iterator.data).Len()
}
