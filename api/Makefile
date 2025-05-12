BUF=buf

.PHONY: proto proto-auth proto-gateway proto-messaging clean


proto:
	$(BUF) generate $(PROTO_REPO)

proto-auth:
	$(BUF) generate --template ./proto/auth/buf.gen.yaml --path proto/auth/auth.proto


proto-gateway:
	cd ./proto/gateway && $(BUF) generate --template buf.gen.yaml $(PROTO_REPO)

proto-messaging:
	cd ./proto/messaging && $(BUF) generate --template buf.gen.yaml $(PROTO_REPO)


clean:
	rm -rf services/**/internal/gen/*