package usecase

import "github.com/leehai1107/The-journal/service/blog/repository"

type IExampleUsecase interface{}

type exampleUsecase struct {
	repo repository.IExampleRepo
}

func NewExampleUsecase(
	repo repository.IExampleRepo,
) IExampleUsecase {
	return &exampleUsecase{
		repo: repo,
	}
}
