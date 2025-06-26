package main

import (
	"fmt"
	"log"
	"sourdough/internal/database"
	"sourdough/internal/handlers"
	"sourdough/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/revrost/go-openrouter"
	"github.com/shareed2k/goth_fiber"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	viper.SetDefault("DEV_MODE", false)

	if viper.GetBool("DEV_MODE") {
		viper.SetDefault("PORT", "3000")
	}

	viper.SetDefault("DB_PATH", "./recipes.db")

	dbPath := viper.GetString("DB_PATH")

	db, err := database.New(dbPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	useProviders()

	// see here for more on what this does: https://github.com/gofiber/storage/blob/main/sqlite3/README.md
	// ...and here for more on why we configure this way: https://docs.giber.io/api/middleware/session
	sessionStore := session.New(session.Config{
		Storage: sqlite3.New(sqlite3.Config{
			Database: dbPath,
		}),
	})

	goth_fiber.SessionStore = sessionStore

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			log.Printf("Error: %v", err)
			return c.Status(code).SendString(err.Error())
		},
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	userRepo := repositories.NewUserRepository(db)
	recipiesRepo := repositories.NewRecipesRepository(db)

	model := viper.GetString("OPENROUTER_MODEL")
	apiKey := viper.GetString("OPENROUTER_API_KEY")
	appName := viper.GetString("OPENROUTER_APP_NAME")

	if model == "" || apiKey == "" || appName == "" {
		log.Fatal("No OpenRouter LLM configuration found. Set OPENROUTER_MODEL, OPENROUTER_API_KEY, and OPENROUTER_APP_NAME environment variables.")
	}

	openRouterClient := openrouter.NewClient(apiKey, openrouter.WithXTitle(appName))

	recipesHandler := handlers.NewRecipesHandler(recipiesRepo, openRouterClient, model)
	authHandler := handlers.NewAuthHandler(userRepo, sessionStore)

	app.Get("/", authHandler.RequireAuth, recipesHandler.GetAllRecipes)
	app.Get("/recipes/:id", recipesHandler.GetRecipe)
	app.Post("/recipes", authHandler.RequireAuth, recipesHandler.PostRecipe)

	app.Get("/login", authHandler.LoginPage)
	app.Get("/auth/:provider", authHandler.Login)
	app.Get("/auth/:provider/callback", authHandler.Callback)
	app.Get("/logout", authHandler.Logout)

	app.Static("/static", "./static")

	port := viper.GetString("PORT")

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func useProviders() {
	googleClientID := viper.GetString("GOOGLE_CLIENT_ID")
	googleClientSecret := viper.GetString("GOOGLE_CLIENT_SECRET")
	baseURL := viper.GetString("BASE_URL")

	if baseURL == "" {
		log.Fatal("No base URL configured. Set BASE_URL environment variable.")
	}

	if viper.GetBool("DEV_MODE") {
		baseURL = fmt.Sprintf("%s:%s", baseURL, viper.GetString("PORT"))
	}

	var providers []goth.Provider

	if googleClientID != "" && googleClientSecret != "" {
		providers = append(providers, google.New(
			googleClientID,
			googleClientSecret,
			baseURL+"/auth/google/callback",
		))
	}

	if len(providers) == 0 {
		log.Fatal("No OAuth providers configured. Set GOOGLE_CLIENT_ID/GOOGLE_CLIENT_SECRET environment variables.")
	}

	goth.UseProviders(providers...)
}
