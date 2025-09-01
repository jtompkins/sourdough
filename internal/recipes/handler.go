package recipes

import (
	"errors"
	"fmt"
	"sourdough/internal/shared"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	repo       *Repository
	llmService *LLMService
}

func NewHandler(repo *Repository, llmService *LLMService) *Handler {
	return &Handler{
		repo:       repo,
		llmService: llmService,
	}
}

func (h *Handler) GetRecipe(c *fiber.Ctx) error {
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

	recipe, err := h.repo.Get(id)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	} else if recipe == nil {
		return c.Status(404).SendString("Recipe not found")
	}

	if user.Id != recipe.UserID {
		return c.Status(403).SendString("Forbidden")
	}

	c.Set("Content-Type", "text/html")
	component := GetRecipeView(recipe)
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func (h *Handler) EditRecipe(c *fiber.Ctx) error {
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

	recipe, err := h.repo.Get(id)

	if err != nil {
		return c.Status(500).SendString(err.Error())
	} else if recipe == nil {
		return c.Status(404).SendString("Recipe not found")
	}

	if user.Id != recipe.UserID {
		return c.Status(403).SendString("Forbidden")
	}

	c.Set("Content-Type", "text/html")
	component := EditRecipeView(recipe)
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func (h *Handler) GetAllRecipes(c *fiber.Ctx) error {
	user, err := h.getCurrentUserFromSession(c)
	if err != nil {
		return err
	}

	recipes, err := h.repo.GetForUser(user.Id)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "text/html")
	component := GetAllRecipesView(recipes)
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func (h *Handler) CreateRecipe(c *fiber.Ctx) error {
	user, err := h.getCurrentUserFromSession(c)
	if err != nil {
		if errors.Is(err, shared.ErrUnauthorized) {
			return c.Status(401).Redirect("/login")
		} else if errors.Is(err, shared.ErrUserNotFound) {
			return c.Status(404).SendString("User not found")
		}
	}

	text := c.FormValue("recipe")

	llmRecipe, err := h.llmService.FormatRecipe(text)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	recipe := llmRecipe.ToRecipe(user.Id)

	result, err := h.repo.Create(&recipe)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Redirect(fmt.Sprintf("/recipes/%d", result.ID))
}

func (h *Handler) UpdateRecipe(c *fiber.Ctx) error {
	user, err := h.getCurrentUserFromSession(c)
	if err != nil {
		if errors.Is(err, shared.ErrUnauthorized) {
			return c.Status(401).Redirect("/login")
		} else if errors.Is(err, shared.ErrUserNotFound) {
			return c.Status(404).SendString("User not found")
		}
	}

	// Get the recipe ID from URL parameters
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).SendString("Invalid recipe ID")
	}

	// Deserialize form data into FormRecipe struct
	var formRecipe FormRecipe
	if err := c.BodyParser(&formRecipe); err != nil {
		return c.Status(400).SendString("Invalid form data: " + err.Error())
	}

	// Convert FormRecipe to Recipe model
	recipe := formRecipe.ToRecipe(user.Id)
	recipe.ID = id

	// Update the recipe in the database (TODO: implement repo.Update)
	_, err = h.repo.Update(&recipe)
	if err != nil {
		return c.Status(500).SendString("Failed to update recipe: " + err.Error())
	}

	// Check if this is an HTMX request
	if c.Get("HX-Request") == "true" {
		// For HTMX requests, send a redirect header to update the URL and content
		c.Set("HX-Redirect", fmt.Sprintf("/recipes/%d", id))
		return c.SendStatus(200)
	} else {
		// For regular requests, redirect back to the recipe page
		return c.Redirect(fmt.Sprintf("/recipes/%d", id))
	}
}

func (h *Handler) getCurrentUserFromSession(c *fiber.Ctx) (*shared.UserInfo, error) {
	userInterface := c.Locals("user")
	if userInterface == nil {
		return nil, shared.ErrUnauthorized
	}

	user, ok := userInterface.(*shared.UserInfo)
	if !ok {
		return nil, shared.ErrUserNotFound
	}

	return user, nil
}
