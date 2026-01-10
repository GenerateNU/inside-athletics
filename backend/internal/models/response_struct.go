package models

type ResponseBody[T any] struct {
	Body *T
}
