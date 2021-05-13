package repos

type IRepository interface{}

type Repository struct {
}

func NewRepository() IRepository {
	repository := &Repository{}
	return repository
}
