# Chirpy

Chirpy is a RESTful API server for a microblogging platform, built with Go. It serves as the backend service handling user authentication, chirps (posts), and administrative metrics for client applications.

## Features

- **Authentication & Authorization**: Secure user registration, password hashing with bcrypt, JWT access tokens, and refresh tokens stored in a database.
- **Chirp Management**: Create, view, filter, and delete posts with automatic profanity filtering.
- **Webhooks**: Handles external payment webhooks to upgrade users to "Chirpy Red" status.
- **Metrics**: Server hit counters and health-check endpoints.

## Tech Stack

- **Language**: Go
- **Database**: PostgreSQL
- **SQL Tooling**: Goose (migrations), sqlc (type-safe Go generation)
- **Security**: JWT (`golang-jwt`), Bcrypt

## Prerequisites

- [Go](https://golang.org/dl/) (version 1.22+)
- [PostgreSQL](https://www.postgresql.org/)

## Environment Variables

Create a `.env` file in the root directory:

```env
PORT=8080
DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=your_jwt_secret_here
POLKA_KEY=your_polka_api_key_here
PLATFORM=dev

```
Getting Started

Clone the repository:

git clone https://github.com/ghost59/chirpy.git
cd chirpy

Run database migrations:

goose -dir sql/schema postgres "$DB_URL" up

Start the server:

go run .

API Endpoints Overview
GET /api/healthz - Health check

POST /api/users - Register a new user

POST /api/login - Authenticate and obtain tokens

GET /api/chirps - List chirps (supports query parameters like author_id and sort)

POST /api/chirps - Create a new chirp (requires auth)

DELETE /api/chirps/{chirpID} - Delete a chirp (author only)

POST /api/polka/webhooks - Upgrade user via webhook



