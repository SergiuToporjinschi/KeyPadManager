package types

type StringOrIntFloat interface {
	~string | ~int | ~int32 | ~int64 | ~int16 | ~int8 | ~uint | ~uint32 | ~uint64 | ~uint16 | ~uint8 | ~float32 | ~float64
}

type AnySet[K StringOrIntFloat] map[K]bool

// type AnySet[K StringOrIntFloat, V any] map[K]KeyValuePair[K, V]

func NewAnySet[K StringOrIntFloat]() AnySet[K] {
	return make(map[K]bool)
}

func NewAnySetWithValues[K StringOrIntFloat](values ...K) AnySet[K] {
	b := NewAnySet[K]()
	for _, value := range values {
		b.Add(value)
	}
	return b
}

func (s AnySet[K]) Add(element K) bool {
	if s.Contains(element) {
		return false
	}
	s[element] = true
	return true
}

func (s AnySet[K]) AddAll(elements ...K) int {
	var addedCnt int
	for _, item := range elements {
		if s.Contains(item) {
			continue
		}
		s[item] = true
		addedCnt++
	}
	return addedCnt
}

func (s AnySet[K]) Remove(element K) bool {
	exists := s.Contains(element)
	if !exists {
		return false
	}
	delete(s, element)
	return true
}

func (s AnySet[K]) Get(element K) bool {
	return s[element]
}

func (s AnySet[K]) Keys() []K {
	keys := make([]K, len(s))
	for item := range s {
		keys = append(keys, item)
	}
	return keys
}

func (s AnySet[K]) Contains(element K) bool {
	_, f := s[element]
	return f
}

func (s AnySet[K]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

func (s AnySet[K]) IsEmpty() bool {
	return s.Size() == 0
}

func (s AnySet[K]) Size() int {
	return len(s)
}
