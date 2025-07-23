package auth

import (
	"sourdough/internal/shared"

	"github.com/gofiber/fiber/v2"
)

type Middleware struct {
	handler *Handler
}

func NewMiddleware(handler *Handler) *Middleware {
	return &Middleware{handler: handler}
}

func (m *Middleware) RequireAuth(c *fiber.Ctx) error {
	sess, err := m.handler.store.Get(c)
	if err != nil {
		return c.Status(401).Redirect("/login")
	}

	authenticated := sess.Get("authenticated")
	if authenticated == nil || authenticated != true {
		return c.Status(401).Redirect("/login")
	}

	user, err := m.handler.getCurrentUser(c)
	if err != nil {
		return c.Status(401).Redirect("/login")
	}

	// Convert to a generic struct to avoid import cycles
	userInfo := &shared.UserInfo{
		Id:       user.Id,
		UserId:   user.UserId,
		Provider: user.Provider,  
	}

	c.Locals("user", userInfo)
	return c.Next()
}