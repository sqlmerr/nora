FROM golang:latest AS build

ARG GOOSE_DBSTRING
ENV GOOSE_DBSTRING=${GOOSE_DBSTRING}

WORKDIR /compiler 

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN CGO_ENABLED=0 GOOS=linux go build -o web ./cmd/web/main.go

FROM scratch AS runtime

WORKDIR /app

COPY --from=build /compiler/.env .
COPY --from=build /compiler/web .
COPY --from=build /compiler/migrations ./migrations


EXPOSE 8000

CMD ["./web"]