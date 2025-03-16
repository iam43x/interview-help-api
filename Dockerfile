FROM golang:alpine AS builder
LABEL stage=gobuilder
RUN apk update --no-cache && apk add --no-cache tzdata
WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app cmd/http/main.go

FROM alpine
RUN apk update --no-cache && apk add --no-cache ca-certificates
COPY --from=builder /build/app /app
ENTRYPOINT ["/app"]