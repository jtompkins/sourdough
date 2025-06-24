# Recipe App

A Go-based recipe management application built with Fiber, SQLite, and server-side rendering using Templ templates.

## Features

- **Go + Fiber**: Fast HTTP web framework
- **SQLite Database**: Lightweight, embedded database
- **Templ Templates**: Type-safe HTML templating
- **HTMX + Alpine.js**: Modern frontend interactivity without build steps
- **OpenRouter Ready**: LLM integration capabilities

## Quick Start

1. **Install dependencies**:
   ```bash
   go mod download
   go install github.com/a-h/templ/cmd/templ@latest
   ```

2. **Generate templates**:
   ```bash
   templ generate
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```

4. **Visit**: http://localhost:3000

## Development Commands

```bash
# Generate Templ templates (run after editing .templ files)
templ generate

# Run development server
go run main.go

# Run with custom port
PORT=8080 go run main.go

# Run with custom database path
DB_PATH=./custom.db go run main.go

# Build for production
go build -o sourdough main.go

# Format code
go fmt ./...

# Run tests
go test ./...

# Watch and regenerate templates (if you have templ installed)
templ generate --watch
```

## Project Structure

```
sourdough/
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── recipes.db             # SQLite database (auto-created)
├── internal/
│   ├── database/
│   │   └── database.go    # Database connection & schema
│   ├── handlers/
│   │   └── handlers.go    # HTTP request handlers
│   └── models/
│       └── recipe.go      # Data models (Recipe, Category)
├── templates/
│   ├── layout.templ       # Base HTML layout
│   ├── index.templ        # Home page template
│   └── *.templ_templ.go   # Generated Go files (don't edit)
└── static/                # Static assets (CSS, JS, images)
```

## Database Schema

The app automatically creates these tables:

- **recipes**: id, title, description, ingredients, instructions, prep_time, cook_time, servings, timestamps
- **categories**: id, name, created_at
- **recipe_categories**: Many-to-many relationship table

## Adding New Features

### 1. **New Routes/Handlers**
Add to `internal/handlers/handlers.go`:
```go
func (h *Handler) NewFeature(c *fiber.Ctx) error {
    // Your handler logic
    return c.SendString("New feature")
}
```

Register in `main.go`:
```go
app.Get("/new-feature", h.NewFeature)
```

### 2. **New Templates**
Create `templates/new-page.templ`:
```templ
package templates

templ NewPage() {
    @Layout("New Page") {
        <h1>New Page Content</h1>
    }
}
```

Run `templ generate` to compile.

### 3. **Database Operations**
Add methods to `internal/handlers/handlers.go`:
```go
func (h *Handler) createRecipe(recipe models.Recipe) error {
    query := `INSERT INTO recipes (title, description, ...) VALUES (?, ?, ...)`
    _, err := h.db.Exec(query, recipe.Title, recipe.Description, ...)
    return err
}
```

### 4. **API Endpoints**
For JSON APIs, add handlers that return JSON:
```go
func (h *Handler) APIGetRecipes(c *fiber.Ctx) error {
    recipes, err := h.getRecentRecipes(10)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(recipes)
}
```

### 5. **OpenRouter/LLM Integration**
Set environment variables:
```bash
export OPENAI_API_KEY="your-openrouter-key"
export OPENAI_BASE_URL="https://openrouter.ai/api/v1"
```

Add LLM service in `internal/`:
```go
// internal/llm/service.go
package llm

import "github.com/sashabaranov/go-openai"

func NewClient() *openai.Client {
    config := openai.DefaultConfig("your-api-key")
    config.BaseURL = "https://openrouter.ai/api/v1"
    return openai.NewClientWithConfig(config)
}
```

## Environment Variables

- `PORT`: Server port (default: 3000)
- `DB_PATH`: Database file path (default: ./recipes.db)
- `OPENAI_API_KEY`: OpenRouter API key
- `OPENAI_BASE_URL`: OpenRouter endpoint URL

## Common Extension Points

1. **Authentication**: Add middleware in `main.go`
2. **Recipe CRUD Operations**: Extend handlers for create/update/delete
3. **Image Upload**: Add file handling for recipe photos
4. **Search**: Implement full-text search in SQLite
5. **Categories Management**: Add category CRUD operations
6. **Recipe Import**: Parse recipes from URLs or text
7. **AI Features**: Generate recipes, suggest substitutions
8. **Export**: Generate PDFs or shopping lists

## Frontend Notes

- **HTMX**: Add `hx-*` attributes for dynamic behavior
- **Alpine.js**: Use `x-data`, `x-show`, etc. for client-side reactivity
- **No Build Process**: All frontend dependencies loaded via CDN
- **Styling**: Add CSS in `templates/layout.templ` or link external stylesheets

## Production Deployment

```bash
# Build binary
go build -o sourdough main.go

# Run with production settings
DB_PATH=/data/recipes.db PORT=80 ./sourdough
```

## Contributing

1. Edit `.templ` files for HTML changes
2. Run `templ generate` after template changes
3. Follow Go conventions for backend code
4. Test locally before committing