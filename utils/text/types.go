package text

type intSlice interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
