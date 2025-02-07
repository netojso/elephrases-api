FROM golang:1.23 AS build

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o api cmd/main.go

# stage imagem final
FROM scratch

WORKDIR /app

COPY --from=build /app/api ./

COPY --from=build /app/.env ./

COPY --from=build /app/database/migrations ./database/migrations

CMD [ "./api" ]