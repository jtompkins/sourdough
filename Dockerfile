# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Install build dependencies: git, nodejs, npm
RUN apk add --no-cache build-base git nodejs npm

WORKDIR /app

# Copy package manager files
COPY go.mod go.sum ./
# Download Go modules
RUN go mod download
# Install templ for code generation
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy package.json and package-lock.json
COPY package.json package-lock.json ./
# Install npm dependencies (for tailwind)
RUN npm ci

# Copy Makefile so we can run make commands
COPY Makefile ./

# Copy the rest of the application source code
COPY . .

# Generate templ files and tailwind css, then build the Go binary
RUN make generate
RUN GOOS=linux make build

# Stage 2: Final
FROM gcr.io/distroless/static-debian11

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/sourdough .

# Copy static assets and templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

# Expose the port the app runs on
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["/app/sourdough"]
