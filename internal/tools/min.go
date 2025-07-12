package tools

func Min[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](a T, b T) T {
	if a < b {
		return a
	} else {
		return b
	}
}
