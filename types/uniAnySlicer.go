package types

type PairKeyValue[K StringOrIntFloat, V any] interface {
	Key() K
	Value() V
}

type UniAnySlice[T StringOrIntFloat, V any] []PairKeyValue[T, V]

func NewUniAnySlice[T StringOrIntFloat, V any]() UniAnySlice[T, V] {
	return []PairKeyValue[T, V]{}
}

func NewUniAnySliceWithValues[T StringOrIntFloat, V any](values ...PairKeyValue[T, V]) UniAnySlice[T, V] {
	return append([]PairKeyValue[T, V]{}, values...)
}

func (s *UniAnySlice[T, V]) Add(element PairKeyValue[T, V]) {
	for _, item := range *s {
		if item.Key() == element.Key() {
			return
		}
	}
	*s = append(*s, element)
}

func (s *UniAnySlice[T, V]) AddAll(elements ...PairKeyValue[T, V]) {
	for _, item := range *s {
		for _, element := range elements {
			if item.Key() == element.Key() {
				continue
			}
		}
	}
	*s = append(*s, elements...)
}

func (s *UniAnySlice[T, V]) RemoveByIndex(index int) {
	if index < 0 || index >= len(*s) {
		return
	}
	*s = append((*s)[:index], (*s)[index+1:]...)
}

func (s *UniAnySlice[T, V]) Remove(element PairKeyValue[T, V]) {
	for i, item := range *s {
		if item.Key() == element.Key() {
			s.RemoveByIndex(i)
			return
		}
	}
}

func (s *UniAnySlice[T, V]) Size() int {
	return len(*s)
}

func (s *UniAnySlice[T, V]) IsEmpty() bool {
	return s.Size() == 0
}

func (s *UniAnySlice[T, V]) Clear() {
	*s = []PairKeyValue[T, V]{}
}

func (s *UniAnySlice[T, V]) Contains(element PairKeyValue[T, V]) bool {
	for _, item := range *s {
		if item.Key() == element.Key() {
			return true
		}
	}
	return false
}

func (s *UniAnySlice[T, V]) Values() []V {
	values := []V{}
	for _, value := range *s {
		values = append(values, value.Value())
	}
	return values
}

func (s *UniAnySlice[T, V]) Keys() []T {
	keys := []T{}
	for _, item := range *s {
		keys = append(keys, item.Key())
	}
	return keys
}

func (s *UniAnySlice[T, V]) GetByIndex(index int) PairKeyValue[T, V] {
	return (*s)[index]
}

func (s *UniAnySlice[T, V]) GetByKey(key T) PairKeyValue[T, V] {
	for _, item := range *s {
		if item.Key() == key {
			return item
		}
	}
	return nil
}
