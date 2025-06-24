# Makefile for the sourdough project

# Using .PHONY declares these targets as not being actual files.
# This is a good practice for targets that are commands.
.PHONY: all run build test fmt templ.generate templ.watch tailwind.watch watch

# The default target executed when you just run `make`
all: build

# Run the application
run:
	go run main.go

# Build the application binary
build:
	go build -o sourdough main.go

# Run the tests
test:
	go test ./...

# Format the Go source code
fmt:
	go fmt ./...

# --- Templ Targets ---

# Generate HTML from templ files
templ.generate:
	templ generate

# Watch for changes in .templ files and regenerate
templ.watch:
	templ generate --watch --proxy=http://localhost:3000 --cmd="go run ."

tailwind.watch:
	npx --yes tailwindcss -i ./static/input.css -o ./static/output.css --watch

watch:
	@echo "Starting watchers in parallel..."
	@trap 'kill %1 %2 2>/dev/null; exit' INT; \
	templ generate --watch --proxy=http://localhost:3000 --cmd="go run ." & \
	npx --yes tailwindcss -i ./static/input.css -o ./static/output.css --watch & \
	wait
