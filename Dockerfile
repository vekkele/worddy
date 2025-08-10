FROM golang:1.24.6-bookworm AS builder

ENV CGO_ENABLED=0

WORKDIR /app

ARG TAILWINDCSS_OS_ARCH
ENV TAILWINDCSS_OS_ARCH=${TAILWINDCSS_OS_ARCH}

COPY go.mod go.sum ./
RUN go mod download

COPY Makefile ./
RUN make install

COPY . ./
RUN make build-prod

FROM debian:bookworm-slim

WORKDIR /app

ENV WORDDY_DB_DSN=""

COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/bin/main ./bin/main
COPY --from=builder /app/internal/store/postgres/migrations ./migrations

EXPOSE 8080

CMD ["sh", "-c", "goose -dir migrations postgres \"$WORDDY_DB_DSN\" up && exec ./bin/main"]
