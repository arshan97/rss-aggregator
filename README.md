# rss-aggregator

A lightweight RSS feed aggregator written in Go. Users register, follow RSS feeds, and the service continuously collects new posts from all followed feeds in the background — so you check one place instead of twenty blogs.

Built to learn Go's approach to backend services: concurrency with goroutines, type-safe SQL with sqlc, and the standard `net/http` stack with chi.

## Architecture

![Architecture diagram](docs/architecture.png)

Two independent parts share one PostgreSQL database:

1. **HTTP API** (chi router) — users, feeds, feed follows, and posts. Protected routes use API key authentication via middleware.
2. **Background scraper** — a long-running goroutine that ticks every minute, picks the feeds that haven't been fetched recently, and fans out to 10 concurrent worker goroutines. Each worker fetches the feed's XML, parses it, and stores new posts (duplicates are skipped).

## Tech stack

- **Go** — chi router, net/http, goroutines + WaitGroup for the scraper
- **PostgreSQL** — data store
- **sqlc** — write raw SQL, generate type-safe Go query code
- **goose** — database migrations

## API endpoints

All routes are under `/v1`. Auth = requires `Authorization: ApiKey <key>` header.

| Method | Endpoint | Auth | Description |
|--------|---------------------------|------|--------------------------------------|
| GET    | /health                   | No   | Readiness check                      |
| POST   | /users                    | No   | Create a user (returns an API key)   |
| GET    | /users                    | Yes  | Get the current user                 |
| POST   | /feeds                    | Yes  | Add a new RSS feed                   |
| GET    | /feeds                    | No   | List all feeds                       |
| POST   | /feed_follows             | Yes  | Follow a feed                        |
| GET    | /feed_follows             | Yes  | List feeds you follow                |
| DELETE | /feed_follows/{id}        | Yes  | Unfollow a feed                      |
| GET    | /posts                    | Yes  | Latest posts from feeds you follow   |

## Getting started

### Prerequisites

- Go 1.21+
- PostgreSQL running locally
- [goose](https://github.com/pressly/goose) and [sqlc](https://sqlc.dev) installed

### Setup

1. Clone the repo and create a `.env` file in the root:

```
PORT=8080
DB_URL=postgres://user:password@localhost:5432/rssagg?sslmode=disable
```

2. Run the database migrations:

```bash
cd sql/schema
goose postgres "postgres://user:password@localhost:5432/rssagg" up
```

3. Start the server:

```bash
go build && ./rss-aggregator
```

The scraper starts automatically alongside the API.

### Try it

```bash
# Create a user (save the api_key from the response)
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "john"}'

# Add a feed
curl -X POST http://localhost:8080/v1/feeds \
  -H "Authorization: ApiKey <your-key>" \
  -H "Content-Type: application/json" \
  -d '{"name": "John Blog", "url": "<your-blog-feed-xml-link>"}'

# Wait a minute for the scraper, then read your posts
curl http://localhost:8080/v1/posts \
  -H "Authorization: ApiKey <your-key>"
```

## License

MIT