# Sourdough

Sourdough is an app that helps users manage recipes. It uses an LLM to clean up recipes that are submitted by the user - many times, recipes are sourced from Instagram or other social media platforms and are imprecise or filled with emoji and other unnecessary information.

## Features

Sourdough is a very simple app. It provides three main features:

* An "all recipes" view, accessible at the root of the app ("/"), which lists all recipes the user has entered and provides a form element to accept new recipes.
* A "recipe" view, accessible at "/recipes/:id", which displays a single recipe and includes a button to edit the recipe.
* An "edit recipe" view, accessible at "/recipes/:id/edit", which allows the user to edit a recipe.

Additionally, the "recipe" view includes a stylesheet with custom print-specific media queries to allow the view to be printed in a more readable format.

## Tech Stack

* Go
* Fiber
* Templ
* Goth
* SQLite
* SQLx
* HTMX
* Alpine.js

Styles are provided using vanilla CSS.

## Architecture

Sourdough uses a common "vertical" structure, where each capability is implemented as a separate package:

* **auth** - Provides authentication using Goth. User login is provided by third-party providers. This package also includes a middleware to check if the user is authenticated before accessing certain routes. The login view is also included here.
* **recipes** - All recipe functionality is implemented here. This package contains models, handlers, the LLM service, and a repository for storing and retrieving recipes. All of the application's views for recipes are here as well.
* **database** - This package contains the database connection functionality and some custom SQLx types for handling JSON in the database.
* **shared** - This package contains shared functionality. It's mostly error types.

## Operations

Most of the application's operations are accessed via the Makefile. Important make targets include:

* `make run` - This runs the application.
* `make build` - This builds the application.
* `make generate` - This runs Templ to generate templates.
* `make db.reset` - This removes the SQLite database. The database is recreated when the app runs if it doesn't exist. IMPORTANT: this operation is destructive and results in data loss. DO NOT run this command unless you are explicitly told to do so.
* `make watch` - This starts the application using Templ's watch mode, allowing the views to be live-reloaded.

### Docker-specific make targets

* `make docker.build` - Builds the docker image.
* `make docker.run` - Runs the docker image once it's built.
* `build.docker` - This is a Docker-specific build command that does Linux cross-compilation for the Linux-based Docker image. My development machine is a Mac, so I need the docker and local builds to be separate targets.

## TODOs

* One day, I'd like to figure out how to add tests.
* I want to implement a feature that allows users to paste in images of recipes and then pass those images to the LLM service for analysis.
* I need to rework the styles for the application to make it look better.
