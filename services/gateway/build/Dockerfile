FROM golang:1.24 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/gateway ./cmd


FROM scratch AS final

WORKDIR /

COPY --from=build /bin/gateway /gateway

EXPOSE 8080

ENTRYPOINT ["/gateway"]