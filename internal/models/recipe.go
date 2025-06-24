package models

import (
	"sourdough/internal/database"
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
