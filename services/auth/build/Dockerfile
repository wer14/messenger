FROM golang:1.24 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/auth ./cmd
RUN go install github.com/grpc-ecosystem/grpc-health-probe@latest

FROM alpine:3.21 AS final

WORKDIR /

RUN apk add --no-cache ca-certificates

COPY --from=build /bin/auth /auth
COPY --from=build /go/bin/grpc-health-probe /usr/local/bin/grpc-health-probe

EXPOSE 8080
EXPOSE 8082

ENTRYPOINT ["/auth"]