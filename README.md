# GProc - Go Process Manager

A lightweight process manager for Go applications, similar to PM2 for Node.js.

## Features

- Start, stop, and restart processes
- Auto-restart failed processes
- Real-time log viewing
- Process monitoring and health checks
- CLI interface

## Installation

```bash
go build -o gproc cmd/main.go
```

## Usage

### Start a process
```bash
./gproc start myapp ./myapp --port 8080
./gproc start myapp ./myapp --auto-restart --max-restarts 10
```

### List processes
```bash
./gproc list
```

### View logs
```bash
./gproc logs myapp
./gproc logs myapp --lines 50
```

### Stop a process
```bash
./gproc stop myapp
```

### Restart a process
```bash
./gproc restart myapp
```

## Project Structure

```
GProc/
├── cmd/
│   └── main.go              # CLI application entry point
├── internal/
│   ├── process/
│   │   └── manager.go       # Process management logic
│   └── logger/
│       └── tail.go          # Log tailing functionality
├── pkg/
│   └── types/
│       └── process.go       # Core data structures
├── logs/                    # Process log files
├── go.mod
└── README.md
```