# Solana Monitor

Real-time Solana wallet balance monitoring with alert rules and Telegram/Email notifications.

## Quick Start

```bash
# Docker (recommended)
cp .env.example .env
# Edit .env with your TELEGRAM_BOT_TOKEN and RPC_URL
docker-compose up -d

# Or manually
cd backend && cp ../.env . && go run cmd/server/main.go
```

## Features

- Monitor multiple Solana wallet addresses
- Configurable alert rules (balance change, threshold, percentage)
- Telegram & Email notifications
- Web dashboard
- Docker deployment

## Configuration

| Variable | Description |
|----------|-------------|
| `RPC_URL` | Solana RPC endpoint |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token |
| `SMTP_*` | Email settings (optional) |
| `HTTP_PROXY` | Proxy for RPC calls (optional) |

## API

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/addresses` | GET/POST | Wallet addresses |
| `/api/rules` | GET/POST | Alert rules |
| `/api/notifications` | GET/POST | Notification channels |
| `/api/alerts` | GET | Alert history |
| `/api/stats` | GET | Dashboard stats |

## Tech Stack

- Backend: Go + Gin + SQLite
- Frontend: Vue 3 + TailwindCSS
- Deployment: Docker + Nginx