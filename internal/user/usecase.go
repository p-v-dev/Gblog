package user

type Usecase interface {
	Execute(input CreateUserInput) (*UserOutput, error)
}

type usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) Execute(input CreateUserInput) (*UserOutput, error) {
	// Lógica de negócio aqui
	return nil, nil
}
