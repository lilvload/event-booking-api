# Event Booking API

A Go-powered REST API for booking events. Users can browse events, create their own, and register for events. Protected routes use JWT authentication.

## Features

- Create, view, update and delete events
- User signup and login with JWT authentication
- Register for and cancel registration from events
- Only the creator of an event can update or delete it

## API Endpoints

| Method | Endpoint                 | Description                | Auth                  |
|--------|--------------------------|----------------------------|-----------------------|
| GET    | `/events`                | Get a list of events       | No                    |
| GET    | `/events/<id>`           | Get a single event         | No                    |
| POST   | `/events`                | Create a new event         | Required              |
| PUT    | `/events/<id>`           | Update an event            | Required (creator)    |
| DELETE | `/events/<id>`           | Delete an event            | Required (creator)    |
| POST   | `/signup`                | Create a new user          | No                    |
| POST   | `/login`                 | Authenticate user (JWT)    | No                    |
| POST   | `/events/<id>/register`  | Register user for event    | Required              |
| DELETE | `/events/<id>/register`  | Cancel registration        | Required              |

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.21 or newer

### Run locally

```bash
git clone https://github.com/your-username/event-booking-api.git
cd event-booking-api
go mod tidy
go run .
```

The server starts on `http://localhost:8080`.

## Authentication

After logging in via `/login`, send the returned token in the `Authorization` header for protected routes:

```
Authorization: <your-jwt-token>
```