# go-profile-api

Profile API written in Go using the Gin framework. Serves a small profile at `/me`, a health check, and a welcome message. This repository was prepared to be pushed to GitHub under `brainox/go-profile-api`.

## Quickstart

Prerequisites:

- Go 1.24+ (toolchain may be declared in `go.mod`)

Run locally:

```bash
go run ./app.go
```

Endpoints:

- `GET /` - welcome message
- `GET /me` - returns profile information and a cat fact
- `GET /health` - health check

Environment variables (optional):

- `USER_EMAIL` - email to return in the profile
- `USER_NAME` - display name
- `USER_STACK` - tech stack string
- `CAT_FACT_API` - external cat fact API URL
- `PORT` - port to listen on (default `8080`)

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.
