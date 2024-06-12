package utility

type IntSet map[int]bool

func NewIntSet() IntSet {
	return make(IntSet)
}

func NewIntSetWithValues(values ...int) IntSet {
	b := make(IntSet)
	for _, value := range values {
		b.Add(value)
	}
	return b
}

func (s IntSet) Add(element int) bool {
	exists := s.Contains(element)
	s[element] = true
	return !exists
}

func (s IntSet) AddAll(elements ...int) {
	for _, item := range elements {
		s[item] = true
	}
}

func (s IntSet) Remove(element int) bool {
	exists := s.Contains(element)
	delete(s, element)
	return exists
}

func (s IntSet) Keys() []int {
	keys := make([]int, 0, len(s))
	for item := range s {
		keys = append(keys, item)
	}
	return keys
}

func (s IntSet) Contains(element int) bool {
	return s[element]
}

func (s IntSet) Values() []int {
	values := make([]int, 0, len(s))
	for value := range s {
		values = append(values, value)
	}
	return values
}
