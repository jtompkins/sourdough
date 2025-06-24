package repositories

import (
	"database/sql"
	"errors"
	"sourdough/internal/database"
	"sourdough/internal/models"
)

type RecipesRepository struct {
	db *database.DB
}

func NewRecipesRepository(db *database.DB) *RecipesRepository {
	return &RecipesRepository{db: db}
}

func (repo *RecipesRepository) Get(id int) (*models.Recipe, error) {
	var recipe models.Recipe

	err := repo.db.Get(&recipe, "SELECT * FROM recipes WHERE id = ?", id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &recipe, nil
}

func (repo *RecipesRepository) GetForUser(userID int) ([]*models.Recipe, error) {
	var recipes []*models.Recipe

	err := repo.db.Select(&recipes, "SELECT * FROM recipes WHERE user_id = ?", userID)

	if err != nil {
		return nil, err
	}

	return recipes, nil
}

func (repo *RecipesRepository) Create(userID int, recipe *models.LLMRecipe) (*models.Recipe, error) {
	result, err := repo.db.Exec("INSERT INTO recipes (user_id, title, ingredients, number_of_ingredients, directions, prep_time, cook_time, servings) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		userID, recipe.Title, recipe.Ingredients, len(recipe.Ingredients), recipe.Directions, recipe.PrepTime, recipe.CookTime, recipe.Servings)

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
