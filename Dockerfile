FROM golang:1.21 AS build

WORKDIR /app

COPY . .
RUN go mod download

COPY ./config /config
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server

FROM alpine:3.18.3

RUN apk --no-cache add ca-certificates

WORKDIR /
COPY --from=build /bin/server /bin/server
COPY --from=build /config /config

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]