# GProc Deployment Guide

## Prerequisites

1. Install [Fly.io CLI](https://fly.io/docs/hands-on/install-flyctl/)
2. Create a Fly.io account and login: `flyctl auth login`

## Deploy Backend

```bash
# From project root
chmod +x deploy-backend.sh
./deploy-backend.sh
```

Or manually:
```bash
flyctl launch --config fly.toml --name gproc-backend
flyctl deploy
```

## Deploy Frontend

```bash
# From fn/ directory
cd fn
chmod +x deploy-frontend.sh
./deploy-frontend.sh
```

Or manually:
```bash
cd fn
flyctl launch --config fly.toml --name gproc-frontend
flyctl deploy
```

## Environment Variables

### Frontend (.env.production)
- `VITE_API_URL`: Backend API URL (https://gproc-backend.fly.dev)
- `VITE_WS_URL`: WebSocket URL (wss://gproc-backend.fly.dev)

### Backend (fly.toml)
- `PORT`: Server port (8080)
- `GIN_MODE`: Gin framework mode (release)

## URLs After Deployment

- **Frontend**: https://gproc-frontend.fly.dev
- **Backend API**: https://gproc-backend.fly.dev/api/v1
- **WebSocket**: wss://gproc-backend.fly.dev/api/v1/ws

## Local Development

### Backend
```bash
go run cmd/main.go cmd/daemon.go daemon
```

### Frontend
```bash
cd fn
npm run dev
```

## Troubleshooting

1. **CORS Issues**: Backend includes CORS middleware for frontend domain
2. **WebSocket Connection**: Ensure WSS is used in production
3. **API Calls**: Check network tab for 404/500 errors
4. **Authentication**: Verify JWT tokens are being sent correctly