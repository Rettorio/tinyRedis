# tinyRedis

Minimal Redis-like TCP server in Go. Single `main` package, no external deps.
8 Harry Potter entries pre-seeded on startup.

## Usage

    go run .                     # start server
    redis-cli -p 6739            # connect

## Commands

| Command | Args | Example | Response |
|---------|------|---------|----------|
| PING   | — | `PING` | `+PONG\r\n` |
| SET    | key value | `SET HP 1 "Sorcerer's Stone"` | `+OK\r\n` |
| GET    | key | `GET "HP 1"` | `$21\r\nSorcerer's Stone\r\n` or `$-1\r\n` if missing |
| DEL    | key | `DEL "HP 1"` | `:1\r\n` (deleted) or `:0\r\n` (not found) |

## Notes

- RESP protocol only — inline/text commands not supported.
- Empty bulk strings (`$0\r\n\r\n`) are rejected.
