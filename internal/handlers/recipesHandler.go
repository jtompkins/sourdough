package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sourdough/internal/models"
	"sourdough/internal/repositories"
	"sourdough/templates"
	"strconv"

	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

const LLM_SYSTEM_PROMPT = `
	You are a helpful model that specializes in formatting recipe text into a structured JSON output.
	When you are given text input, if it looks like a recipe, you will do the following steps:
		1. Clean up the formatting of individual ingredients, normalizing the measurements to American standards
		2. Simplify individual steps in the instructions where it makes sense, but DO NOT remove or skip steps
		3. If you cannot determine a value for any of fields, output an empty string ("") for the value, DO NOT substitute any other value or skip the field
		4. Return your modified version of the recipe in JSON format, adhering to the following schema:
			{
				"title": "string",
				"prepTime": "string", # in hours and minutes
				"cookTime": "string", # in hours and minutes
				"servings": "number",
				"ingredients": [
					"string"
				],
				"instructions": [
					"string"
				]
			}
`

type RecipesHandler struct {
	recipesRepo  *repositories.RecipesRepository
	openAIClient *openai.Client
	llmModel     string
}

func NewRecipesHandler(recipesRepo *repositories.RecipesRepository, openAIClient *openai.Client, llmModel string) *RecipesHandler {
	return &RecipesHandler{
		recipesRepo:  recipesRepo,
		openAIClient: openAIClient,
		llmModel:     llmModel,
	}
}

func (h *RecipesHandler) GetRecipe(c *fiber.Ctx) error {
	user, err := h.getCurrentUserFromSession(c)
	if err != nil {
		return err
	}

	idParam := c.Params("id")

	if idParam == "" {
		return c.Status(400).SendString("Missing recipe ID")
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).SendString("Invalid recipe ID")
	}

	recipe, err := h.recipesRepo.Get(id)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	} else if recipe == nil {
		return c.Status(404).SendString("Recipe not found")
	}

	if user.Id != recipe.UserID {
		return c.Status(403).SendString("Forbidden")
	}

	c.Set("Content-Type", "text/html")
	component := templates.Recipe(recipe)
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func (h *RecipesHandler) GetAllRecipes(c *fiber.Ctx) error {
	user, err := h.getCurrentUserFromSession(c)
	if err != nil {
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
	var llmRecipe models.LLMRecipe

	schema, err := jsonschema.GenerateSchemaForType(llmRecipe)
	if err != nil {
		return models.LLMRecipe{}, err
	}

	req := openai.ChatCompletionRequest{
		Model: h.llmModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: LLM_SYSTEM_PROMPT,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: recipe,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "recipe",
				Schema: schema,
				Strict: true,
			},
		},
	}

	resp, err := h.openAIClient.CreateChatCompletion(
		context.Background(),
		req,
	)

	if err != nil {
		return models.LLMRecipe{}, err
	}

	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &llmRecipe)

	if err != nil {
		return models.LLMRecipe{}, err
	}

	return llmRecipe, nil
}
