package models

type LLMRecipe struct {
	Title       string   `json:"title"`
	Ingredients []string `json:"ingredients"`
	Directions  []string `json:"directions"`
	PrepTime    string   `json:"prepTime"`
	CookTime    string   `json:"cookTime"`
	Servings    int      `json:"servings"`
}
