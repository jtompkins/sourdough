package repositories

import (
	"database/sql"
	"errors"
	"sourdough/internal/database"
	"sourdough/internal/models"
)

type UsersRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) Get(id int) (*models.User, error) {
	var user models.User

	err := r.db.Get(&user, "SELECT * FROM users WHERE id = ?", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (repo *UsersRepository) GetByProviderId(userId string) (*models.User, error) {
	var user models.User

	err := repo.db.Get(&user, "SELECT * from users where user_id = ?", userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (repo *UsersRepository) Create(userId string, provider string) (*models.User, error) {
	query := "INSERT INTO users (user_id, provider) VALUES (?, ?)"
	tx, err := repo.db.Beginx()
	if err != nil {
		return nil, err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	result, err := tx.Exec(query, userId, provider)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return repo.Get(int(id))
}
