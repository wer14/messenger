PROTO_REPO=../../proto

# Путь к buf CLI, если требуется локальный вызов
BUF=buf

.PHONY: proto proto-auth proto-gateway proto-messaging clean

## Общая генерация (если понадобится для всего)
proto:
	$(BUF) generate $(PROTO_REPO)

## Генерация только для auth-service
proto-auth:
	cd ./proto/auth && $(BUF) generate --template buf.gen.yaml $(PROTO_REPO)

## Генерация только для gateway (включает grpc-gateway)
proto-gateway:
	cd ./proto/gateway && $(BUF) generate --template buf.gen.yaml $(PROTO_REPO)

## Генерация только для messaging-service
proto-messaging:
	cd ./proto/messaging && $(BUF) generate --template buf.gen.yaml $(PROTO_REPO)

## Очистка сгенерированных файлов
clean:
	rm -rf services/**/internal/gen/*