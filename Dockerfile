# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Install build dependencies: git
RUN apk add --no-cache build-base git

WORKDIR /app

# Copy package manager files
COPY go.mod go.sum ./
# Download Go modules
RUN go mod download
# Install templ for code generation
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy Makefile so we can run make commands
COPY Makefile ./

# Copy the rest of the application source code
COPY . .

# Generate templ files and build the Go binary
RUN make docker-build

# Stage 2: Final
FROM gcr.io/distroless/static-debian11

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/sourdough .

# Copy static assets and templates


# Expose the port the app runs on
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["/app/sourdough"]
