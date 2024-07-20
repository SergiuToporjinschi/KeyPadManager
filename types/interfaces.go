package types

type StringOrIntFloat interface {
	~string | ~int | ~int32 | ~int64 | ~int16 | ~int8 | ~uint | ~uint32 | ~uint64 | ~uint16 | ~uint8 | ~float32 | ~float64
}

type Displayable interface {
	GetDisplayText() string
}

type Uniqueable[T StringOrIntFloat] interface {
	GetID() T
}
