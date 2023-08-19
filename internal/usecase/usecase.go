package usecase

type UseCase struct {
	storage storage
	logger  logger.Logger
}

type storage interface {
}

func New() UseCase {
	return UseCase{}
}
