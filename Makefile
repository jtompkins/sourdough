# Makefile for the sourdough project

# Using .PHONY declares these targets as not being actual files.
# This is a good practice for targets that are commands.
.PHONY: all run build generate test fmt templ.generate templ.watch tailwind.generate tailwind.watch watch

# The default target executed when you just run `make`
all: build

# Run the application
run: templ.generate tailwind.generate
	go run main.go

build: templ.generate tailwind.generate
	go build -o sourdough main.go

generate: templ.generate tailwind.generate

# Run the tests
test:
	go test ./...

# Format the Go source code
fmt:
	go fmt ./...

db.reset:
	rm *.db

templ.generate:
	templ generate

templ.watch:
	templ generate --watch --proxy=http://localhost:3000 --cmd="go run ."

tailwind.generate:
	npx --yes tailwindcss -i ./static/input.css -o ./static/output.css

tailwind.watch:
	npx --yes tailwindcss -i ./static/input.css -o ./static/output.css --watch

watch:
	@echo "Starting watchers in parallel..."
	@trap 'kill %1 %2 2>/dev/null; exit' INT; \
	templ generate --watch --proxy=http://localhost:3000 --cmd="go run ." & \
	npx --yes tailwindcss -i ./static/input.css -o ./static/output.css --watch & \
	wait
