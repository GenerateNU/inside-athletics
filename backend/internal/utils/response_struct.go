package utils

type ResponseBody[T any] struct {
	Body *T
}
