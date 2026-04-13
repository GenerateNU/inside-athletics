package utils

func MapList[T any, Y any](original []T, apply func(T) Y) []Y {
	mappedList := make([]Y, 0)
	for _, value := range original {
		mappedList = append(mappedList, apply(value))
	}
	return mappedList
}
