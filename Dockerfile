FROM golang:1.26.5-alpine AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . ./
RUN go run ./cmd/phraseforge-db --database-file /out/phraseforge.db
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/phraseforge-api ./cmd/phraseforge-api

FROM alpine:3.22

WORKDIR /app

COPY --from=build --chown=nobody:nobody /out/phraseforge-api /app/phraseforge-api
COPY --from=build --chown=nobody:nobody /out/phraseforge.db /app/data/phraseforge.db

USER nobody:nobody

EXPOSE 8080

ENTRYPOINT ["/app/phraseforge-api"]
