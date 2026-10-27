FROM golang:1.26.4-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/app ./cmd/app

FROM alpine:3.20
RUN adduser -D app
USER app
COPY --from=build /out/app /app
EXPOSE 8080
CMD ["/app"]