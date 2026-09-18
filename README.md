# LivePoll

Real-time polling application.

## Tech Stack

- React (frontend, Vite)
- Go + Gin (backend)
- MongoDB (persistent storage)
- Redis (fast counters + Pub/Sub for live updates)
- WebSocket (push live results to open browsers)

## Features

- Signup / Login (JWT auth)
- Create poll (with backend validation)
- Share poll link
- Vote
- Live results — no refresh needed
- WebSocket + Redis Pub/Sub

## Project structure

```
LivePoll/
├── frontend/   React app (Vite)
├── backend/    Go + Gin API
└── README.md
```

## How to run

### 1. Set up MongoDB Atlas and Redis Cloud

Create free databases on MongoDB Atlas and Redis Cloud, then grab their connection strings.

### 2. Backend

```
cd backend
cp .env.example .env      # fill in MONGODB_URI, REDIS_URL, JWT_SECRET
go mod tidy
go run main.go
```

Backend runs on `http://localhost:8080`.

### 3. Frontend

```
cd frontend
cp .env.example .env      # defaults already point at localhost:8080
npm install
npm run dev
```

Frontend runs on `http://localhost:5173`.

## How live updates work

1. A vote hits `POST /api/polls/:id/vote`.
2. The backend increments the count in **Redis** (instant) and publishes the
   new count on a Redis Pub/Sub channel (`poll:<id>`).
3. Every open results page is connected over **WebSocket** to
   `GET /api/ws/polls/:id`, which is subscribed to that same channel.
4. The backend forwards the Pub/Sub message straight to the WebSocket, and
   the React page updates the bar chart — no refresh needed.
5. In the background, the same vote is also persisted to **MongoDB** so
   counts survive a restart.

## API summary

| Method | Route                    | Auth | Description                     |
|--------|--------------------------|------|----------------------------------|
| POST   | /api/signup              | No   | Create an account                |
| POST   | /api/login                | No   | Get a JWT                        |
| POST   | /api/polls                | Yes  | Create a poll                    |
| GET    | /api/polls                | Yes  | List your polls                  |
| GET    | /api/polls/:id            | No   | Get one poll                     |
| POST   | /api/polls/:id/vote        | No   | Cast a vote                      |
| GET    | /api/ws/polls/:id          | No   | WebSocket for live results       |
