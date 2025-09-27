package recipes

import (
	"database/sql"
	"errors"
	"sourdough/internal/database"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) Get(id int) (*Recipe, error) {
	var recipe Recipe

	err := repo.db.Get(&recipe, "SELECT * FROM recipes WHERE id = ?", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &recipe, nil
}

func (repo *Repository) Delete(id int) (bool, error) {
	result, err := repo.db.Exec("DELETE FROM recipes WHERE id = ?", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (repo *Repository) GetForUser(userID int) ([]*Recipe, error) {
	var recipes []*Recipe

	err := repo.db.Select(&recipes, "SELECT * FROM recipes WHERE user_id = ?", userID)

	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (repo *Repository) Search(userID int, searchTerm string) ([]*Recipe, error) {
	var recipes []*Recipe

	likeClause := "%" + searchTerm + "%"

	err := repo.db.Select(&recipes, "SELECT * FROM recipes WHERE user_id = ? and title LIKE ?", userID, likeClause)

	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (repo *Repository) Create(recipe *Recipe) (*Recipe, error) {

	// Use SQLx's NamedExec to automatically map struct fields to query parameters
	result, err := repo.db.NamedExec(
		"INSERT INTO recipes (user_id, title, ingredients, number_of_ingredients, directions, prep_time, cook_time, servings) VALUES (:user_id, :title, :ingredients, :number_of_ingredients, :directions, :prep_time, :cook_time, :servings)",
		recipe,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Fetch and return the inserted recipe
	return repo.Get(int(id))
}

func (repo *Repository) Update(recipe *Recipe) (*Recipe, error) {
	// Use SQLx's NamedExec to automatically map struct fields to query parameters
	_, err := repo.db.NamedExec(
		"UPDATE recipes SET title = :title, ingredients = :ingredients, number_of_ingredients = :number_of_ingredients, directions = :directions, prep_time = :prep_time, cook_time = :cook_time, servings = :servings WHERE id = :id",
		recipe,
	)
	if err != nil {
		return nil, err
	}

	// Fetch and return the inserted recipe
	return repo.Get(recipe.ID)
}
