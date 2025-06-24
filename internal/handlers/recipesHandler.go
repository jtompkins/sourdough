package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sourdough/internal/models"
	"sourdough/internal/repositories"
	"sourdough/templates"

	"github.com/gofiber/fiber/v2"
	"github.com/revrost/go-openrouter"
	"github.com/revrost/go-openrouter/jsonschema"
)

const LLM_SYSTEM_PROMPT = `

`

type RecipesHandler struct {
	recipesRepo      *repositories.RecipesRepository
	openRouterClient *openrouter.Client
	llmModel         string
	responseSchema   *jsonschema.Definition
}

func NewRecipesHandler(recipesRepo *repositories.RecipesRepository, openRouterClient *openrouter.Client, llmModel string) *RecipesHandler {
	var schemaType models.LLMRecipe
	responseSchema, err := jsonschema.GenerateSchemaForType(schemaType)

	if err != nil {
		log.Fatalf("GenerateSchemaForType error: %v", err)
	}

	return &RecipesHandler{
		recipesRepo:      recipesRepo,
		openRouterClient: openRouterClient,
		llmModel:         llmModel,
		responseSchema:   responseSchema,
	}
}

func (h *RecipesHandler) GetRecipe(c *fiber.Ctx) error {
	var sampleRecipe = models.Recipe{
		Title:    "Chocolate-Covered Pickle Ice Cream with Mustard Sprinkles",
		PrepTime: "47 minutes",
		CookTime: "15 minutes",
		Ingredients: []string{
			"2 cups vanilla ice cream, melted and confused",
			"1½ cups dill pickle juice, chilled",
			"12 large pickles, diced into perfect cubes",
			"8 oz dark chocolate, melted backwards",
			"½ cup yellow mustard, frozen into tiny pearls",
			"3 tablespoons pickle brine reduction",
			"1 teaspoon vanilla extract (the sad kind)",
			"¼ cup crushed pretzel confusion",
			"2 drops green food coloring (optional, for extra pickle vibes)",
		},
		Directions: []string{
			"Gently whisper to the melted ice cream until it remembers how to be cold again (approximately 15 minutes).",
			"Stir in pickle juice using only counterclockwise motions while humming your favorite pickle song.",
			"Fold in the diced pickles as if you're tucking them into bed for a long winter's nap.",
			"Drizzle the backwards-melted chocolate in zigzag patterns that spell out \"why\" in cursive.",
			"Freeze the mixture while standing on one foot and thinking about regrets (20 minutes).",
			"Using a melon baller, scoop into serving bowls and sprinkle with frozen mustard pearls.",
			"Garnish with crushed pretzel confusion and serve immediately to unsuspecting guests.",
			"Document their facial expressions for science.",
		},
		NumberOfIngredients: 9,
		Servings:            6,
	}

	c.Set("Content-Type", "text/html")
	component := templates.Recipe(&sampleRecipe)
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func (h *RecipesHandler) GetAllRecipes(c *fiber.Ctx) error {
	user, err := h.getCurrentUserFromSession(c)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return c.Status(401).Redirect("/login")
		} else if errors.Is(err, ErrUserNotFound) {
			return c.Status(404).SendString("User not found")
		}

		return err
	}

	recipes, err := h.recipesRepo.GetForUser(user.Id)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "text/html")
	component := templates.MyRecipes(recipes)
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func (h *RecipesHandler) PostRecipe(c *fiber.Ctx) error {
	user, err := h.getCurrentUserFromSession(c)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			return c.Status(401).Redirect("/login")
		} else if errors.Is(err, ErrUserNotFound) {
			return c.Status(404).SendString("User not found")
		}
	}

	text := c.FormValue("recipe")

	recipe, err := h.formatWithLLM(text)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	result, err := h.recipesRepo.Create(user.Id, &recipe)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Redirect(fmt.Sprintf("/recipes/%d", result.ID))
}

func (h *RecipesHandler) getCurrentUserFromSession(c *fiber.Ctx) (*models.User, error) {
	userInterface := c.Locals("user")
	if userInterface == nil {
		return nil, ErrUnauthorized
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (h *RecipesHandler) formatWithLLM(recipe string) (models.LLMRecipe, error) {
	req := openrouter.ChatCompletionRequest{
		Model: h.llmModel,
		Messages: []openrouter.ChatCompletionMessage{
			{
				Role:    openrouter.ChatMessageRoleSystem,
				Content: openrouter.Content{Text: LLM_SYSTEM_PROMPT},
			},
			{
				Role:    openrouter.ChatMessageRoleUser,
				Content: openrouter.Content{Text: recipe},
			},
		},
		ResponseFormat: &openrouter.ChatCompletionResponseFormat{
			Type: openrouter.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openrouter.ChatCompletionResponseFormatJSONSchema{
				Name:   "recipe",
				Schema: h.responseSchema,
				Strict: true,
			},
		},
	}

	resp, err := h.openRouterClient.CreateChatCompletion(
		context.Background(),
		req,
	)

	if err != nil {
		return models.LLMRecipe{}, err
	}

	var llmRecipe models.LLMRecipe

	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content.Text), &llmRecipe)

	if err != nil {
		return models.LLMRecipe{}, err
	}

	return llmRecipe, nil
}
