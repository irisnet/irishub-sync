package helper

type Iterator interface {
	HasNext() bool
	Next() interface{}
	Remove()
	Get() []interface{}
}

func newIntIterator(array *[]int) Iterator {
	var d []interface{}
	for _, data := range *array {
		d = append(d, data)
	}
	return &ArrayIterator{
		data: &d,
	}
}

func newStringIterator(array *[]string) Iterator {
	var d []interface{}
	for _, data := range *array {
		d = append(d, data)
	}
	return &ArrayIterator{
		data: &d,
	}
}

func GetIterator(array interface{}) Iterator {
	switch array.(type) {
	case *[]int:
		return newIntIterator(array.(*[]int))
	case *[]string:
		return newStringIterator(array.(*[]string))

	}
	return nil
}

type ArrayIterator struct {
	index int
	data  *[]interface{}
}

func (iterator *ArrayIterator) HasNext() bool {
	return iterator.index < len(*iterator.data)
}

func (iterator *ArrayIterator) Next() (v interface{}) {
	v = (*iterator.data)[iterator.index]
	iterator.index++
	return v
}

func (iterator *ArrayIterator) Remove() {
	var data = *iterator.data
	if iterator.index == 0 {
		data = data[1:]
		iterator.data = &data
		return
	} else if iterator.index == len(data) {
		data = data[0 : len(data)-1]
		iterator.data = &data
		return
	}
	frontData, afterData := data[0:iterator.index-1], data[iterator.index:]
	newData := append(frontData, afterData...)
	iterator.data = &newData
}

func (iterator *ArrayIterator) Get() []interface{} {
	return *iterator.data
}
