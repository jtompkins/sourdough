# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based recipe application using:
- **Fiber** web framework for REST API
- **SQLite** database with go-sqlite3 driver
- **Templ** for server-side HTML templating
- **HTMX** and **Alpine.js** for frontend interactivity (via CDN)

## Development Commands

```bash
# Generate Templ & Tailwind templates
make generate

# Run the application
make run

# Run with custom database path
DB_PATH=./custom.db go run main.go

# Run with custom port
PORT=8080 go run main.go

# Build the application
go build -o sourdough main.go

# Format code
go fmt ./...

# Run tests
go test ./...
```

## Project Structure

```
├── cmd/                    # Application entrypoints
├── internal/
│   ├── database/          # Database connection and setup
│   ├── handlers/          # HTTP handlers
│   └── models/           # Data models
├── templates/            # Templ templates (.templ files)
├── static/              # Static assets
├── main.go             # Application entry point
└── recipes.db          # SQLite database (created automatically)
```

## Architecture Notes

- **Database**: SQLite with recipes, categories, and recipe_categories tables
- **Templates**: Use Templ for type-safe HTML generation
- **Frontend**: HTMX for dynamic interactions, Alpine.js for client-side reactivity
- **No build steps**: Frontend uses CDN resources, no bundling required
- **OpenRouter**: Ready for LLM integration via sashabaranov/go-openai client
