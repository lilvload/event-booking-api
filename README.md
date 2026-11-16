# Event Booking API

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.12-008ECF)
![SQLite](https://img.shields.io/badge/SQLite-embedded-003B57?logo=sqlite&logoColor=white)
![JWT](https://img.shields.io/badge/auth-JWT-black?logo=jsonwebtokens)

A REST API for booking events, written in **Go** with the **Gin** framework and an embedded **SQLite** database.

Anyone can browse events. Registered users can create their own events, edit or delete the events they created, and sign up for (or cancel their registration for) any event. Protected routes are secured with JWT, and passwords are stored as bcrypt hashes.

## Table of contents

- [Features](#features)
- [Tech stack](#tech-stack)
- [Project structure](#project-structure)
- [Getting started](#getting-started)
- [Configuration](#configuration)
- [Authentication](#authentication)
- [API reference](#api-reference)
- [Usage example](#usage-example)
- [Database schema](#database-schema)
- [Testing the API](#testing-the-api)
- [Known limitations and roadmap](#known-limitations-and-roadmap)
- [Author](#author)

## Features

- Full CRUD for events: list, get by ID, create, update, delete
- User signup and login
- Passwords hashed with bcrypt before they are stored
- JWT authentication (HS256, tokens valid for 2 hours) implemented as a Gin middleware
- Ownership check: only the user who created an event can update or delete it
- Registration for an event and cancellation of a registration
- SQLite database that is created automatically on first start, with no external services to install
- Request validation: all event and user fields are required
- Ready-to-run `.http` request files for every endpoint

## Tech stack

| Area | Tool |
|------|------|
| Language | [Go](https://go.dev/) 1.26 |
| HTTP framework | [Gin](https://github.com/gin-gonic/gin) v1.12 |
| Database | SQLite via [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) (pure Go, no CGO required) and `database/sql` |
| Authentication | [golang-jwt/jwt](https://github.com/golang-jwt/jwt) v5 |
| Password hashing | [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) (`golang.org/x/crypto`) |

## Project structure

```
event-booking-api/
├── main.go              entry point: initializes the database and starts the server
├── db/
│   └── db.go            SQLite connection, connection pool settings, table creation
├── models/
│   ├── event.go         Event model: save, list, get by ID, update, delete, register, cancel
│   └── user.go          User model: save with hashed password, credential validation
├── routes/
│   ├── routes.go        route table and the authenticated route group
│   ├── events.go        event handlers
│   ├── users.go         signup and login handlers
│   └── register.go      event registration handlers
├── middlewares/
│   └── auth.go          JWT authentication middleware
├── utils/
│   ├── hash.go          bcrypt helpers
│   └── jwt.go           token generation and verification
├── api-test/            sample HTTP requests for every endpoint
├── go.mod
└── go.sum
```

## Getting started

### Prerequisites

- [Go](https://go.dev/dl/) 1.26 or newer (the version declared in `go.mod`)

No database server and no C compiler are needed: the SQLite driver is pure Go.

### Run locally

```bash
git clone https://github.com/lilvload/event-booking-api.git
cd event-booking-api
go mod download
go run .
```

The server starts on `http://localhost:8080`. On first start it creates the `api.db` SQLite file in the working directory together with all tables.

### Build a binary

```bash
go build -o bin/event-booking-api .
./bin/event-booking-api
```

To start from an empty database, stop the server and delete `api.db`.

## Configuration

The project currently has no environment variables. All settings are constants in the code:

| Setting | Value | Location |
|---------|-------|----------|
| Port | `8080` | `main.go` |
| Database file | `api.db` (working directory) | `db/db.go` |
| Connection pool | 10 open / 5 idle connections | `db/db.go` |
| JWT signing key | hard-coded constant | `utils/jwt.go` |
| Token lifetime | 2 hours | `utils/jwt.go` |
| bcrypt cost | 14 | `utils/hash.go` |

Gin runs in debug mode by default. Set `GIN_MODE=release` to turn off debug logging.


## Authentication

1. Create an account with `POST /signup`.
2. Log in with `POST /login` to receive a JWT.
3. Send the token in the `Authorization` header of every protected request.

```
Authorization: <your-jwt-token>
```

The header must contain the raw token **without** the `Bearer ` prefix. A missing, invalid or expired token results in `401 Unauthorized`.

## API reference

Base URL: `http://localhost:8080`

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/events` | List all events | No |
| `GET` | `/events/:id` | Get a single event | No |
| `POST` | `/events` | Create an event | Required |
| `PUT` | `/events/:id` | Update an event | Required (creator only) |
| `DELETE` | `/events/:id` | Delete an event | Required (creator only) |
| `POST` | `/events/:id/register` | Register the current user for an event | Required |
| `DELETE` | `/events/:id/register` | Cancel the current user's registration | Required |
| `POST` | `/signup` | Create a user | No |
| `POST` | `/login` | Log in and receive a JWT | No |

Request bodies are JSON. Field names in request bodies are case-insensitive (`name` and `Name` both work); responses use capitalized field names.

### Users

#### `POST /signup`

```json
{
  "email": "user@example.com",
  "password": "test123"
}
```

| Status | Body |
|--------|------|
| `201 Created` | `{"message": "User created successfully!"}` |
| `400 Bad Request` | A field is missing or the JSON is invalid |
| `500 Internal Server Error` | The user could not be saved, for example the email is already taken |

#### `POST /login`

Takes the same body as `/signup`.

```json
{
  "message": "Login succesful!",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

| Status | Body |
|--------|------|
| `200 OK` | Message and token |
| `400 Bad Request` | A field is missing or the JSON is invalid |
| `401 Unauthorized` | `{"message": "Credentials invalid"}` |

### Events

#### Event object

```json
{
  "ID": 1,
  "Name": "Go meetup",
  "Description": "Talks",
  "Location": "Kyiv",
  "DateTime": "2026-12-01T18:00:00Z",
  "UserID": 1
}
```

`UserID` is the ID of the user who created the event. It is taken from the token, not from the request body.

#### `GET /events`

Returns `200 OK` with an array of event objects (`null` when there are no events yet).

#### `GET /events/:id`

| Status | Body |
|--------|------|
| `200 OK` | Event object |
| `400 Bad Request` | The ID is not a number |
| `500 Internal Server Error` | The event could not be fetched, including when it does not exist |

#### `POST /events`

All four fields are required. `dateTime` must be an RFC 3339 timestamp.

```json
{
  "name": "Go meetup",
  "description": "Talks",
  "location": "Kyiv",
  "dateTime": "2026-12-01T18:00:00.000Z"
}
```

| Status | Body |
|--------|------|
| `201 Created` | `{"message": "Event created!", "event": { ... }}` |
| `400 Bad Request` | A field is missing or the JSON is invalid |
| `401 Unauthorized` | No valid token |

#### `PUT /events/:id`

Takes the same body as `POST /events`; all fields are required, so this is a full replacement.

| Status | Body |
|--------|------|
| `200 OK` | `{"message": "Event updated successfully!"}` |
| `400 Bad Request` | Invalid ID or request body |
| `401 Unauthorized` | No valid token, or the current user is not the creator of the event |
| `500 Internal Server Error` | The event could not be fetched (including when it does not exist) or updated |

#### `DELETE /events/:id`

| Status | Body |
|--------|------|
| `200 OK` | `{"message": "Event deleted successfully!"}` |
| `400 Bad Request` | Invalid ID |
| `401 Unauthorized` | No valid token, or the current user is not the creator of the event |
| `500 Internal Server Error` | The event could not be fetched, including when it does not exist |

### Registrations

#### `POST /events/:id/register`

No request body. Registers the user identified by the token.

| Status | Body |
|--------|------|
| `201 Created` | `{"message": "Registered!"}` |
| `400 Bad Request` | Invalid ID |
| `401 Unauthorized` | No valid token |
| `500 Internal Server Error` | The event could not be fetched, including when it does not exist |

#### `DELETE /events/:id/register`

No request body.

| Status | Body |
|--------|------|
| `200 OK` | `{"message": "Cancelled!"}` |
| `400 Bad Request` | Invalid ID |
| `401 Unauthorized` | No valid token |

All error responses have the form `{"message": "..."}`.


## Database schema

Three tables are created automatically in `api.db`:

| Table | Columns |
|-------|---------|
| `users` | `id` (PK, autoincrement), `email` (unique, not null), `password` (bcrypt hash, not null) |
| `events` | `id` (PK, autoincrement), `name`, `description`, `location`, `dateTime` (all not null), `user_id` → `users.id` |
| `registrations` | `id` (PK, autoincrement), `event_id` → `events.id`, `user_id` → `users.id` |

```
users 1 ──── * events
users 1 ──── * registrations * ──── 1 events
```

## Testing the API

The [`api-test/`](api-test) folder contains one `.http` file per endpoint. They can be run with the [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) extension for VS Code or the built-in HTTP client of GoLand / IntelliJ IDEA.

| File | Request |
|------|---------|
| `create-user.http` | `POST /signup` |
| `login.http` | `POST /login` |
| `get-events.http` | `GET /events` |
| `get-single-event.http` | `GET /events/:id` |
| `create-event.http` | `POST /events` |
| `update-event.http` | `PUT /events/:id` |
| `delete-event.http` | `DELETE /events/:id` |
| `register.http` | `POST /events/:id/register` |
| `cancel-registration.http` | `DELETE /events/:id/register` |

The tokens saved in these files have expired. Run `login.http` first and paste the new token into the `Authorization` header of the protected requests.

There are no automated tests yet.

## Author

**Vladyslav Bakunets** · [@lilvload](https://github.com/lilvload)
