FROM golang:alpine as builder

WORKDIR /app
COPY ./go.mod .
RUN go mod download
RUN apk add --no-cache gcc musl-dev
COPY . .
RUN CGO_ENABLED=1 GOOS=linux  go build -o api

FROM alpine

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/api /api

EXPOSE 9093

ENTRYPOINT ["/api"]