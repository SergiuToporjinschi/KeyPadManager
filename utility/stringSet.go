package utility

type StringSet map[string]bool

func NewStringSet() StringSet {
	return make(StringSet)
}

func NewStringSetWithValues(values ...string) StringSet {
	b := make(StringSet)
	for _, value := range values {
		b.Add(value)
	}
	return b
}

func (s StringSet) Add(element string) bool {
	exists := s.Contains(element)
	s[element] = true
	return !exists
}

func (s StringSet) AddAll(elements ...string) {
	for _, item := range elements {
		s[item] = true
	}
}

func (s StringSet) Remove(element string) bool {
	exists := s.Contains(element)
	delete(s, element)
	return exists
}

func (s StringSet) Keys() []string {
	keys := make([]string, 0, len(s))
	for item := range s {
		keys = append(keys, item)
	}
	return keys
}

func (s StringSet) Contains(element string) bool {
	return s[element]
}

func (s StringSet) Values() []string {
	values := make([]string, 0, len(s))
	for value := range s {
		values = append(values, value)
	}
	return values
}
