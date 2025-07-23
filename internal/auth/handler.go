package auth

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/markbates/goth"
	"github.com/shareed2k/goth_fiber"
)

type Handler struct {
	userRepo *Repository
	store    *session.Store
}

func NewHandler(repo *Repository, store *session.Store) *Handler {
	return &Handler{userRepo: repo, store: store}
}

func (h *Handler) LoginPage(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/html")
	component := LoginView()
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func (h *Handler) Login(c *fiber.Ctx) error {
	provider := c.Params("provider")
	if provider == "" {
		return c.Status(400).SendString("Provider is required")
	}

	return goth_fiber.BeginAuthHandler(c)
}

func (h *Handler) Callback(c *fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		log.Printf("Auth callback error: %v", err)
		return c.Status(500).Redirect("/login?error=auth_failed")
	}

	dbUser, err := h.findOrCreateUser(user)
	if err != nil {
		log.Printf("Database error: %v", err)
		return c.Status(500).Redirect("/login?error=db_error")
	}

	sess, err := h.store.Get(c)
	if err != nil {
		log.Printf("Session error: %v", err)
		return c.Status(500).Redirect("/login?error=session_error")
	}

	sess.Set("user_id", dbUser.Id)
	sess.Set("authenticated", true)

	if err := sess.Save(); err != nil {
		log.Printf("Session save error: %v", err)
		return c.Status(500).Redirect("/login?error=session_save")
	}

	return c.Redirect("/")
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	if err := goth_fiber.Logout(c); err != nil {
		log.Printf("Logout error: %v", err)
	}

	sess, err := h.store.Get(c)
	if err != nil {
		log.Printf("Session error during logout: %v", err)
		return c.Redirect("/")
	}

	sess.Delete("user_id")
	sess.Delete("authenticated")
	sess.Destroy()

	if err := sess.Save(); err != nil {
		log.Printf("Session save error during logout: %v", err)
	}

	return c.Redirect("/")
}

func (h *Handler) findOrCreateUser(gothUser goth.User) (*User, error) {
	userId := gothUser.Provider + ":" + gothUser.UserID

	user, err := h.userRepo.GetByProviderId(userId)

	if err != nil {
		return nil, err
	} else if user != nil {
		return user, nil
	}

	user, err = h.userRepo.Create(userId, gothUser.Provider)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *Handler) getCurrentUser(c *fiber.Ctx) (*User, error) {
	sess, err := h.store.Get(c)
	if err != nil {
		return nil, err
	}

	userIDInterface := sess.Get("user_id")
	if userIDInterface == nil {
		return nil, fiber.NewError(401, "Not authenticated")
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		return nil, fiber.NewError(401, "Invalid user session")
	}

	user, err := h.userRepo.Get(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}