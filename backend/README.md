# LMS Backend
 
A learning management system (LMS) backend written in Go. Teachers create courses, admins approve them, and students enroll.
 
**Stack:** Go, Gin, PostgreSQL, JWT authentication, Docker Compose.
 
## Features
 
- Sign up, log in, refresh tokens and log out (JWT access token + refresh token cookie)
- Three roles: `student` (default), `teacher` and `admin`
- Teachers create courses; new courses are `pending` until an admin approves them
- Students enroll in published courses (no duplicate enrollments)
- Students see their own enrollments; the course's teacher and admins see the course roster
## Project structure
 
```
backend/
├── cmd/app/            # main.go, config (environment variables)
├── docker/             # Dockerfile and docker-compose.yml
├── internal/
│   ├── user/           # users, auth, refresh tokens
│   ├── course/         # courses
│   └── enrollment/     # enrollments
│       ├── domain/         # entities, repository interfaces, services
│       ├── infrastructure/ # PostgreSQL repositories
│       └── interfaces/http # Gin handlers
├── migrations/         # SQL migrations (golang-migrate)
└── pkg/                # errorhandler, jwt, middleware
```
 
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
 
The API is now available at `http://localhost:8080`.
 
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
 
Never commit `.env`. It is listed in `.gitignore`.
 
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
| GET | `/courses` | anyone | List published courses |
| GET | `/courses/:id` | anyone | Get one course |
| POST | `/courses` | teacher | Create a course (starts as `pending`). Body: `title`, `description` |
| GET | `/courses/mine` | teacher | List your own courses |
| GET | `/courses/pending` | admin | List courses waiting for approval |
| PATCH | `/courses/:id/approve` | admin | Publish a pending course |
| DELETE | `/courses/:id` | teacher (owner) or admin | Delete a course |
 
### Enrollments
 
| Method | Path | Who | Description |
| --- | --- | --- | --- |
| POST | `/courses/:id/enroll` | student | Enroll in a published course (409 if already enrolled) |
| GET | `/enrollments/mine` | logged in | List your enrollments |
| GET | `/courses/:id/enrollments` | teacher (owner) or admin | List the students enrolled in a course |
 
### Errors
 
Errors have this shape:
 
```json
{ "code": 400, "message": "...", "status_code": 400 }
```
 
`code` identifies the error (for example `300` course not found, `400` already enrolled) and `status_code` is the HTTP status.
 
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
4. The student enrolls (`POST /courses/:id/enroll`); enrolling twice returns 409.
5. The student lists their courses (`GET /enrollments/mine`); the teacher lists the roster (`GET /courses/:id/enrollments`).
## Tests
 
```bash
cd backend
go test ./...
```
 
The service layer has unit tests that use fake repositories, so they need no database.
 
## Database migrations
 
SQL migrations live in `backend/migrations` and run automatically with Docker Compose. To add one:
 
```bash
migrate create -ext sql -dir migrations -seq create_something_table
```
 
Never edit a migration that has already been applied; add a new one instead.
 