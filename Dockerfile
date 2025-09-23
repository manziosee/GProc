FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy only necessary files for backend
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o gproc cmd/main.go cmd/daemon.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/gproc .

EXPOSE 8080
CMD ["./gproc", "daemon"]