FROM golang:1.26-alpine3.23 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o tic-tac-toe

FROM alpine:3.23

WORKDIR /app

RUN adduser -D tic-tac-toe && \
    chown -R tic-tac-toe:tic-tac-toe /app

RUN apk add --no-cache ca-certificates tzdata

# Copy the compiled binary and static assets from the build stage
COPY --from=build --chown=tic-tac-toe:tic-tac-toe /app/tic-tac-toe .
COPY --from=build --chown=tic-tac-toe:tic-tac-toe /app/database ./database
COPY --from=build --chown=tic-tac-toe:tic-tac-toe /app/assets ./assets

# Switch to non-root user
USER tic-tac-toe

EXPOSE 8000
