package recipes

import (
	"sourdough/internal/database"
	"strings"
	"time"
)

type Recipe struct {
	ID                  int                        `db:"id"`
	UserID              int                        `db:"user_id"`
	Title               string                     `db:"title"`
	Ingredients         database.JSONArray[string] `db:"ingredients"`
	NumberOfIngredients int                        `db:"number_of_ingredients"`
	Directions          database.JSONArray[string] `db:"directions"`
	PrepTime            string                     `db:"prep_time"`
	CookTime            string                     `db:"cook_time"`
	Servings            int                        `db:"servings"`
	CreatedAt           time.Time                  `db:"created_at"`
	UpdatedAt           time.Time                  `db:"updated_at"`
}

type FormRecipe struct {
	Title               string `form:"title"`
	Ingredients         string `form:"ingredients"`
	NumberOfIngredients int    `form:"number_of_ingredients"`
	Directions          string `form:"directions"`
	PrepTime            string `form:"prep_time"`
	CookTime            string `form:"cook_time"`
	Servings            int    `form:"servings"`
}

func (r FormRecipe) ToRecipe(userID int) Recipe {
	return Recipe{
		UserID:              userID,
		Title:               r.Title,
		Ingredients:         database.JSONArray[string](strings.Split(r.Ingredients, "\n")),
		NumberOfIngredients: r.NumberOfIngredients,
		Directions:          database.JSONArray[string](strings.Split(r.Directions, "\n")),
		PrepTime:            r.PrepTime,
		CookTime:            r.CookTime,
		Servings:            r.Servings,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

type LLMRecipe struct {
	Title       string   `json:"title"`
	Ingredients []string `json:"ingredients"`
	Directions  []string `json:"directions"`
	PrepTime    string   `json:"prepTime"`
	CookTime    string   `json:"cookTime"`
	Servings    int      `json:"servings"`
}

func (r LLMRecipe) ToRecipe(userID int) Recipe {
	return Recipe{
		UserID:              userID,
		Title:               r.Title,
		Ingredients:         database.JSONArray[string](r.Ingredients),
		NumberOfIngredients: len(r.Ingredients),
		Directions:          database.JSONArray[string](r.Directions),
		PrepTime:            r.PrepTime,
		CookTime:            r.CookTime,
		Servings:            r.Servings,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}
