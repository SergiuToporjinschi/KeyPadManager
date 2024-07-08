package types

type StringOrIntFloat interface {
	~string | ~int | ~int32 | ~int64 | ~int16 | ~int8 | ~uint | ~uint32 | ~uint64 | ~uint16 | ~uint8 | ~float32 | ~float64
}

type KeyValuePair[K StringOrIntFloat, V any] interface {
	Key() K
}

type AnySet[K StringOrIntFloat, V any] map[K]KeyValuePair[K, V]

func NewAnySet[K StringOrIntFloat, V any]() AnySet[K, V] {
	return make(AnySet[K, V])
}

func NewAnySetWithValues[K StringOrIntFloat, V any](values ...KeyValuePair[K, V]) AnySet[K, V] {
	b := make(AnySet[K, V])
	for _, value := range values {
		b.Add(value)
	}
	return b
}

func (s AnySet[K, V]) Add(element KeyValuePair[K, V]) bool {
	exists := s.Contains(element)
	s[element.Key()] = element
	return !exists
}

func (s AnySet[K, V]) AddAll(elements ...KeyValuePair[K, V]) {
	for _, item := range elements {
		s[item.Key()] = item
	}
}

func (s AnySet[K, V]) Remove(element KeyValuePair[K, V]) bool {
	exists := s.Contains(element)
	delete(s, element.Key())
	return exists
}

func (s AnySet[K, V]) Get(key K) V {
	return s[key].(V)
}

func (s AnySet[K, V]) Keys() []K {
	keys := make([]K, 0, len(s))
	for item := range s {
		keys = append(keys, item)
	}
	return keys
}

func (s AnySet[K, V]) Contains(element KeyValuePair[K, V]) bool {
	_, f := s[element.Key()]
	return f
}

func (s AnySet[K, V]) Values() []V {
	values := make([]V, 0, len(s))
	for _, value := range s {
		values = append(values, value.(V))
	}
	return values
}

func (s AnySet[K, V]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

func (s AnySet[K, V]) IsEmpty() bool {
	return s.Size() == 0
}

func (s AnySet[K, V]) Size() int {
	return len(s)
}
