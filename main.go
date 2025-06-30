package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"sourdough/internal/database"
	"sourdough/internal/handlers"
	"sourdough/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	openai "github.com/sashabaranov/go-openai"
	"github.com/shareed2k/goth_fiber"
	"github.com/spf13/viper"
)

//go:embed static
var embededStatic embed.FS

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8080") // Default port for production
	viper.SetDefault("DEV_MODE", false)
	viper.SetDefault("DB_PATH", "./recipes.db")
	viper.SetDefault("LLM_PROVIDER_BASE_URL", "https://openrouter.ai/api/v1")

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

	model := viper.GetString("LLM_PROVIDER_MODEL")
	apiKey := viper.GetString("LLM_PROVIDER_API_KEY")
	apiURL := viper.GetString("LLM_PROVIDER_BASE_URL")

	if model == "" || apiKey == "" || apiURL == "" {
		log.Fatal("No LLM configuration found. Set LLM_PROVIDER_MODEL, LLM_PROVIDER_API_KEY, and LLM_PROVIDER_BASE_URL environment variables.")
	}

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = apiURL
	openAIClient := openai.NewClientWithConfig(config)

	recipesHandler := handlers.NewRecipesHandler(recipiesRepo, openAIClient, model)
	authHandler := handlers.NewAuthHandler(userRepo, sessionStore)

	app.Get("/", authHandler.RequireAuth, recipesHandler.GetAllRecipes)
	app.Get("/recipes/:id", authHandler.RequireAuth, recipesHandler.GetRecipe)
	app.Post("/recipes", authHandler.RequireAuth, recipesHandler.PostRecipe)

	app.Get("/login", authHandler.LoginPage)
	app.Get("/auth/:provider", authHandler.Login)
	app.Get("/auth/:provider/callback", authHandler.Callback)
	app.Get("/logout", authHandler.Logout)

	if viper.GetBool("DEV_MODE") {
		app.Static("/static", "./static")
	} else {
		staticFS, err := fs.Sub(embededStatic, "static")
		if err != nil {
			log.Fatal(err)
		}
		app.Use("/static", filesystem.New(filesystem.Config{
			Root:   http.FS(staticFS),
			Browse: true,
		}))
	}

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
