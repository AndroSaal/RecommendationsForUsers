package repository

type Repository interface{}

type UserRepository struct {
	relDB *PostgresDB
}

func NewUserRepository(db *PostgresDB) *UserRepository {
	return &UserRepository{
		relDB: db,
	}
}
