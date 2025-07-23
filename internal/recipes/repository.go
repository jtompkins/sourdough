package recipes

import (
	"database/sql"
	"encoding/json"
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

func (repo *Repository) GetForUser(userID int) ([]*Recipe, error) {
	var recipes []*Recipe

	err := repo.db.Select(&recipes, "SELECT * FROM recipes WHERE user_id = ?", userID)

	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (repo *Repository) Create(userID int, recipe *LLMRecipe) (*Recipe, error) {
	ingredientsJson, err := json.Marshal(recipe.Ingredients)
	if err != nil {
		return nil, err
	}

	directionsJson, err := json.Marshal(recipe.Directions)
	if err != nil {
		return nil, err
	}

	result, err := repo.db.Exec("INSERT INTO recipes (user_id, title, ingredients, number_of_ingredients, directions, prep_time, cook_time, servings) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		userID, recipe.Title, ingredientsJson, len(recipe.Ingredients), directionsJson, recipe.PrepTime, recipe.CookTime, recipe.Servings)

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