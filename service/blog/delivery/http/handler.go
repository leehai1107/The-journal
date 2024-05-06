package http

import "github.com/leehai1107/The-journal/service/blog/usecase"

type IHandler interface{}

type Handler struct {
	usecase usecase.IExampleUsecase
}

func NewHandler(usecase usecase.IExampleUsecase) IHandler {
	return &Handler{
		usecase: usecase,
	}
}
