FROM golang:1.24-alpine3.22 as builder
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
RUN apk add --no-cache curl
COPY /config ./config
COPY /migrations ./migrations
COPY /cmd ./cmd
COPY /config ./config
COPY /internal ./internal
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /usr/local/bin/app ./cmd
FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /app/migrations/ /migrations
COPY --from=builder /usr/local/bin/app /app
COPY --from=builder /usr/bin/curl /usr/bin/curl
COPY --from=builder /lib/ /lib/
COPY --from=builder /usr/lib/ /usr/lib/
CMD ["/app"]