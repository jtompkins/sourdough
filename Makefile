# Makefile for the sourdough project

# Using .PHONY declares these targets as not being actual files.
# This is a good practice for targets that are commands.
.PHONY: all run build generate test fmt watch build.docker docker.build docker.run

# The default target executed when you just run `make`
all: build

# Run the application
run: generate
	go run main.go

build: generate
	go build -o sourdough main.go

build.docker: generate
	CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o sourdough main.go

# Run the tests
test:
	go test ./...

# Format the Go source code
fmt:
	go fmt ./...

db.reset:
	rm *.db

generate:
	templ generate

watch:
	templ generate --watch --proxy=http://localhost:3000 --cmd="go run ."

docker.build:
	docker build -t sourdough-app .

docker.run: docker.build
	docker run -p 8080:8080 --env-file .env sourdough-app
