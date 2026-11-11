# LMS Backend

A learning management system (LMS) backend written in Go. Teachers create courses and lessons, admins approve courses, and students enroll and read the lessons.

**Stack:** Go, Gin, PostgreSQL, JWT authentication, Docker Compose, GitHub Actions.

## Features

- Sign up, log in, refresh tokens and log out (JWT access token + refresh token cookie)
- Three roles: `student` (default), `teacher` and `admin`
- Teachers create courses; new courses are `pending` until an admin approves them
- Students enroll in published courses (no duplicate enrollments) and can un-enroll
- Students see their own enrollments; the course's teacher and admins see the course roster with the students' names
- Teachers add ordered lessons to their courses; enrolled students, the course's teacher and admins can read them
- Teachers schedule live class sessions for their courses and start/end them; enrolled students, the course's teacher and admins can see them
- Paginated course list (`?page=` and `?limit=`)
- Unit tests for the service layer and a CI/CD pipeline (GitHub Actions)

## Project structure

```
.
├── .github/workflows/      # ci.yml (checks) and cd.yml (publishes the Docker image)
└── backend/
    ├── cmd/app/            # main.go, config (environment variables)
    ├── docker/             # Dockerfile and docker-compose.yml
    ├── internal/
    │   ├── user/           # users, auth, refresh tokens
    │   ├── course/         # courses
    │   ├── enrollment/     # enrollments and the course roster
    │   ├── lesson/         # lessons inside courses
    │   └── session/        # scheduled live classes of a course
    │       ├── application/dto/          # request bodies
    │       ├── domain/                   # entities, repository interfaces, services, response values
    │       ├── infrastructure/postgres/  # PostgreSQL repositories
    │       └── interfaces/http/          # Gin handlers
    ├── migrations/         # SQL migrations (golang-migrate)
    └── pkg/                # errorhandler, jwt, middleware, pagination
```

Every module under `internal/` has the same layers: `domain` (business rules), `infrastructure` (database), `interfaces` (HTTP) and `application` (request bodies).

## Getting started

### Run with Docker (recommended)

Requirements: Docker with the Compose plugin.

```bash
cd backend
cp .env.example .env
```

Edit `.env`: set `DB_PASSWORD` and generate the two JWT secrets:

```bash
openssl rand -base64 48   # run twice: JWT_ACCESS_SECRET and JWT_REFRESH_SECRET
```

Start everything (PostgreSQL, migrations and the app):

```bash
docker compose -f docker/docker-compose.yml --env-file .env up --build
```

The API is now available at `http://localhost:8080`. New migrations are applied automatically every time you start the stack.

Stop it with `Ctrl+C`, or remove the containers with:

```bash
docker compose -f docker/docker-compose.yml down       # keeps the database data
docker compose -f docker/docker-compose.yml down -v    # also deletes the database data
```

### Run without Docker

Requirements: Go (see `go.mod`), PostgreSQL and the [`migrate`](https://github.com/golang-migrate/migrate) CLI.

```bash
cd backend
cp .env.example .env          # set DB_HOST=localhost and your database settings
migrate -path migrations -database "postgres://USER:PASSWORD@localhost:5432/lms?sslmode=disable" up
go run ./cmd/app
```

### Run the published image

Every push to `main` publishes the backend image to GitHub Container Registry:

```bash
docker pull ghcr.io/sinapourmahmoud/lms-project-backend:latest
```

Tags: `latest` (newest `main`) and `sha-<commit>` (one exact commit). The image needs the same environment variables as below and a PostgreSQL database that already has the migrations applied.

## Configuration

All settings are environment variables (loaded from `backend/.env`). See `.env.example`.

| Variable | Default | Description |
| --- | --- | --- |
| `APPLICATION_PORT` | `8080` | Port the API listens on |
| `DB_HOST` | `localhost` | Database host (Docker Compose overrides it with `db`) |
| `DB_PORT` | `5432` | Database port |
| `DB_USER` | `username` | Database user |
| `DB_PASSWORD` | `password` | Database password |
| `DB_NAME` | `lms` | Database name |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode |
| `JWT_ACCESS_SECRET` | none | Secret that signs access tokens (required) |
| `JWT_REFRESH_SECRET` | none | Secret that signs refresh tokens (required) |
| `ACCESS_TOKEN_TTL_MIN` | `15` | Access token lifetime in minutes |
| `REFRESH_TOKEN_TTL_HR` | `168` | Refresh token lifetime in hours |

The app refuses to start if the two JWT secrets are empty. Never commit `.env`. It is listed in `.gitignore`.

## API

Protected endpoints need the header `Authorization: Bearer <access_token>`.

### Auth

| Method | Path | Who | Description |
| --- | --- | --- | --- |
| POST | `/signup` | anyone | Create an account (role `student`). Body: `username`, `email`, `password` (min 8 characters) |
| POST | `/login` | anyone | Returns `access_token` and sets the `refresh_token` cookie. Body: `email`, `password` |
| POST | `/refresh` | anyone with the cookie | Returns a new access token |
| POST | `/logout` | anyone with the cookie | Revokes the refresh token |
| GET | `/whoami` | logged in | Returns your user ID |

### Courses

| Method | Path | Who | Description |
| --- | --- | --- | --- |
| GET | `/courses` | anyone | List published courses, one page at a time (see Pagination below) |
| GET | `/courses/:id` | anyone | Get one course |
| POST | `/courses` | teacher | Create a course (starts as `pending`). Body: `title`, `description` |
| GET | `/courses/mine` | teacher | List your own courses |
| GET | `/courses/pending` | admin | List courses waiting for approval |
| PATCH | `/courses/:id/approve` | admin | Publish a pending course |
| DELETE | `/courses/:id` | teacher (owner) or admin | Delete a course (its lessons and sessions are deleted with it) |

#### Pagination

`GET /courses` takes two optional query parameters:

| Parameter | Default | Description |
| --- | --- | --- |
| `page` | `1` | Page number, starting at 1 |
| `limit` | `10` | Courses per page, from 1 to 100 |

Example: `GET /courses?page=2&limit=5`. The response contains the courses of that page and the paging information:

```json
{
  "items": [
    {
      "id": "6f1c2a1e-0000-0000-0000-000000000000",
      "title": "Learn Go",
      "description": "Go from zero",
      "teacher_id": "9a7b3c2d-0000-0000-0000-000000000000",
      "status": "published"
    }
  ],
  "page": 2,
  "limit": 5,
  "total": 37
}
```

`total` is the number of published courses in all pages. A page past the end returns `items: []`. A `page` or `limit` that is not a valid number, or a `limit` above 100, returns `400`.

### Enrollments

| Method | Path | Who | Description |
| --- | --- | --- | --- |
| POST | `/courses/:id/enroll` | student | Enroll in a published course (409 if already enrolled) |
| DELETE | `/courses/:id/enroll` | student | Un-enroll from a course (404 if not enrolled) |
| GET | `/enrollments/mine` | logged in | List your enrollments |
| GET | `/courses/:id/enrollments` | teacher (owner) or admin | List the students enrolled in a course, with their names |

`GET /courses/:id/enrollments` returns the course roster (an empty course returns `[]`), ordered by enrollment date:

```json
[
  {
    "student_id": "6f1c2a1e-0000-0000-0000-000000000000",
    "username": "ali",
    "email": "ali@example.com",
    "enrolled_at": "2026-10-28T10:00:00Z"
  }
]
```

Another teacher or a student calling this endpoint gets `403`. An unknown course gets `404`.

### Lessons

Lessons belong to a course and are shown in the order of their `position` (1, 2, 3, ...). A course cannot have two lessons with the same position.

| Method | Path | Who | Description |
| --- | --- | --- | --- |
| POST | `/courses/:id/lessons` | teacher (owner) | Add a lesson. Body: `title`, `content`, `position` (min 1) |
| GET | `/courses/:id/lessons` | owner teacher, admin or enrolled student | List the lessons of a course, ordered by `position` |
| GET | `/lessons/:id` | owner teacher, admin or enrolled student | Get one lesson |
| PUT | `/lessons/:id` | teacher (owner) | Replace a lesson. Body: `title`, `content`, `position` (all required) |
| DELETE | `/lessons/:id` | teacher (owner) or admin | Delete a lesson (204 with no body) |

A lesson looks like this:

```json
{
  "id": "0b8f7c4e-0000-0000-0000-000000000000",
  "course_id": "6f1c2a1e-0000-0000-0000-000000000000",
  "title": "Introduction",
  "content": "Welcome to the course.",
  "position": 1,
  "created_at": "2026-11-04T10:00:00Z",
  "updated_at": null
}
```

Responses: `201` for create, `200` for read and update, `204` for delete. A student who is not enrolled, or a teacher who does not own the course, gets `403`. An unknown course or lesson gets `404`. A position that is already used in the course gets `409`.

### Sessions

A session is one scheduled live class of a course. Its status moves in one direction only: `scheduled` -> `live` -> `ended`.

| Method | Path | Who | Description |
| --- | --- | --- | --- |
| POST | `/courses/:id/sessions` | teacher (owner) | Schedule a class. Body: `title`, `starts_at` (RFC 3339, must be in the future) |
| GET | `/courses/:id/sessions` | owner teacher, admin or enrolled student | List the sessions of a course, ordered by `starts_at`, with `?page=` and `?limit=` (same response shape as the course list) |
| GET | `/sessions/:id` | owner teacher, admin or enrolled student | Get one session |
| PUT | `/sessions/:id` | teacher (owner) | Change `title` and `starts_at` (only while `scheduled`) |
| POST | `/sessions/:id/start` | teacher (owner) | Move a `scheduled` session to `live` |
| POST | `/sessions/:id/end` | teacher (owner) | Move a `live` session to `ended` |
| DELETE | `/sessions/:id` | teacher (owner) or admin | Delete a session (204 with no body; not allowed while `live`) |

A session looks like this:

```json
{
  "id": "4a2174e3-0000-0000-0000-000000000000",
  "course_id": "6f1c2a1e-0000-0000-0000-000000000000",
  "title": "Live Q&A",
  "starts_at": "2026-11-20T18:00:00Z",
  "status": "scheduled",
  "started_at": null,
  "ended_at": null,
  "created_at": "2026-11-11T10:00:00Z",
  "updated_at": null
}
```

Responses: `201` for create, `200` for read, update, start and end, `204` for delete. A start time in the past gets `400`. A student who is not enrolled, or a teacher who does not own the course, gets `403`. An unknown course or session gets `404`. An action that is not allowed in the current status (starting a live session, ending a scheduled one, editing or deleting a live one) gets `409`.

### Errors

Errors have this shape:

```json
{ "code": 400, "message": "...", "status_code": 400 }
```

`code` identifies the error and `status_code` is the HTTP status. They are different numbers.

| Code | Meaning | Status |
| --- | --- | --- |
| 1 | Internal system error | 500 |
| 2 | Invalid request body format | 400 |
| 3 | Validation error | 400 |
| 4 | Database error | 500 |
| 100-106 | Auth errors: invalid credentials, email or username already used, invalid or expired tokens | 401 or 409 |
| 107 | Forbidden: you do not have permission | 403 |
| 200 | User not found | 404 |
| 300 | Course not found | 404 |
| 400 | Already enrolled in this course | 409 |
| 404 | Not enrolled in this course | 404 |
| 500 | Lesson not found | 404 |
| 501 | Another lesson in this course already has this position | 409 |
| 600 | Session not found | 404 |
| 601 | This action is not allowed in the current state of the session | 409 |

## Creating the first admin

Everyone who signs up is a `student`, and there is no endpoint to create admins. Sign up normally, then promote the user in the database:

```bash
docker compose -f docker/docker-compose.yml exec db \
  psql -U postgres -d lms -c "UPDATE users SET role = 'admin' WHERE email = 'you@example.com';"
```

Use your own `DB_USER` and `DB_NAME` if they differ, and log in again to get a token with the new role. The same command with `'teacher'` creates a teacher.

## Typical flow

1. Sign up a teacher, a student and an admin, and set their roles as described above.
2. The teacher creates a course (`POST /courses`).
3. The admin approves it (`PATCH /courses/:id/approve`).
4. The teacher adds lessons (`POST /courses/:id/lessons`); a repeated position returns 409.
5. The student enrolls (`POST /courses/:id/enroll`); enrolling twice returns 409.
6. The student reads the lessons (`GET /courses/:id/lessons`); before enrolling this returns 403.
7. The teacher lists the roster with names and emails (`GET /courses/:id/enrollments`).
8. The student can leave the course (`DELETE /courses/:id/enroll`); leaving a course they are not in returns 404.

## Tests

```bash
cd backend
go test ./...
```

The service layer (enrollments, courses, lessons and sessions) has unit tests that use fake repositories, so they need no database. The pagination helper has its own unit tests.

## CI/CD

Two GitHub Actions workflows live in `.github/workflows/`:

| Workflow | Runs on | What it does |
| --- | --- | --- |
| `ci.yml` | pushes to `main` and pull requests to `main` | Checks formatting (`gofmt`), runs `go vet`, `go build` and `go test ./...` (job `test`), then builds the Docker image (job `docker`) |
| `cd.yml` | pushes to `main` | Builds the Docker image and publishes it to `ghcr.io/sinapourmahmoud/lms-project-backend` with the tags `latest` and `sha-<commit>` |

`main` is protected: changes go through a pull request, and the `test` and `docker` checks must pass before it can be merged.

Before opening a pull request, run the same checks locally from `backend/`:

```bash
gofmt -l . && go vet ./... && go build ./... && go test ./...
```

`gofmt -l .` should print nothing.

## Database migrations

SQL migrations live in `backend/migrations` and run automatically with Docker Compose. To add one:

```bash
migrate create -ext sql -dir migrations -seq create_something_table
```

Never edit a migration that has already been applied; add a new one instead.

Current migrations: users, refresh tokens, courses, enrollments, lessons and sessions.