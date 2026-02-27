package repository

type Repository struct {
	User UserRepository
}

func NewRepository(db PgxIface) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
