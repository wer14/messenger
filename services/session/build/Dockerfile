FROM golang:1.24 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/session ./cmd


FROM scratch AS final

WORKDIR /

COPY --from=build /bin/session /session

EXPOSE 8080

ENTRYPOINT ["/session"]