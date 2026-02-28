package repository

type Repository struct {
	User    UserRepository
	Session SessionRepository
}

func NewRepository(db PgxIface) *Repository {
	return &Repository{
		User:    NewUserRepository(db),
		Session: NewSessionRepository(db),
	}
}
