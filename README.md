# Sourdough

A Go-based recipe management application built with Fiber, SQLite, and server-side rendering using Templ templates.

## Runbook

### Critical links

- [Production app](https://sourdoughed.com)
- [Fly console](https://fly.io/apps/sourdough)
- [Google cloud console](https://console.cloud.google.com/auth/clients)
- [Cloudflare console](https://dash.cloudflare.com/e22ddf4182166a6d342828d2b40b4313/sourdoughed.com)
- [Font Awesome dashboard](https://fontawesome.com/kits/994b24a8e7/settings)

### Required environment configuration

- `DEV_MODE`: Set to `true` to enable development mode. This should only ever be set locally.
- `DB_PATH`: The path to your SQLite database file.
- `BASE_URL`: The base URL of your deployed app, as seen in the Fly dashboard or the DNS configuration on Cloudflare.
- `GOOGLE_CLIENT_ID`: The client ID for your Google OAuth app.
- `GOOGLE_CLIENT_SECRET`: The client secret for your Google OAuth app.
- `LLM_PROVIDER_BASE_URL`: The base URL for an OpenAI API-compatible LLM provider.
- `LLM_PROVIDER_API_KEY`: The API key for your LLM provider.
- `LLM_PROVIDER_MODEL`: The model name for your LLM provider.

### Operations

- To build the app locally: `make build`
- To run the app locally with TEMPL generation: `make watch`
- To build the docker image locally: `make docker.build` (Don't confuse this with `build.docker`, which is used by the Fly config to do the required Linux cross-compilation)
- To deploy the app: `fly deploy`
- To set env vars on Fly: `fly secrets set <key>=<value>` (these can be copied straight from your .env file)

## Features

- **Go + Fiber**: Fast HTTP web framework
- **SQLite Database**: Lightweight, embedded database
- **Templ Templates**: Type-safe HTML templating
- **HTMX + Alpine.js**: Modern frontend interactivity without build steps
- **LLM Ready**: LLM integration capabilities

## Quick Start

1. **Install dependencies**:
   ```bash
   go mod download
   go install github.com/a-h/templ/cmd/templ@latest
   ```

2. **Run the application**:
   ```bash
   make run
   ```

3. **Visit**: http://localhost:8080
