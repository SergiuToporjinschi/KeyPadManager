package types

type UniSlice[T StringOrIntFloat] []T

func NewUniSlice[T StringOrIntFloat]() UniSlice[T] {
	return []T{}
}

func NewUniSliceWithValues[T StringOrIntFloat](values ...T) UniSlice[T] {
	return append([]T{}, values...)
}

func (s *UniSlice[T]) Add(element T) {
	for _, item := range *s {
		if item == element {
			return
		}
	}
	*s = append(*s, element)
}

func (s *UniSlice[T]) AddAll(elements ...T) {
	for _, item := range *s {
		for _, element := range elements {
			if item == element {
				return
			}
		}
	}
	*s = append(*s, elements...)
}

func (s *UniSlice[T]) RemoveByIndex(index int) {
	if index < 0 || index >= len(*s) {
		return
	}
	*s = append((*s)[:index], (*s)[index+1:]...)
}

func (s *UniSlice[T]) Remove(element T) {
	for i, item := range *s {
		if item == element {
			s.RemoveByIndex(i)
			return
		}
	}
}

func (s *UniSlice[T]) Get(index int) T {
	return (*s)[index]
}

func (s *UniSlice[T]) Size() int {
	return len(*s)
}
