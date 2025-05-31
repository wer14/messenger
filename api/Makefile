BUF=buf
VENDOR_PROTO_PATH := $(CURDIR)/vendor.protobuf

.PHONY: proto proto-auth proto-gateway proto-messaging clean vendor-googleapis

vendor-googleapis:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/googleapis/googleapis $(VENDOR_PROTO_PATH)/googleapis &&\
	cd $(VENDOR_PROTO_PATH)/googleapis &&\
	git checkout
	mv $(VENDOR_PROTO_PATH)/googleapis/google $(VENDOR_PROTO_PATH)
	rm -rf $(VENDOR_PROTO_PATH)/googleapis

proto:
	$(BUF) generate $(PROTO_REPO)

proto-auth:
	$(BUF) generate --template ./proto/auth/buf.gen.yaml --path proto/auth/auth.proto


proto-gateway:
	$(BUF) generate --template ./proto/gateway/buf.gen.yaml --path proto/gateway/gateway.proto


proto-messaging:
	cd ./proto/messaging && $(BUF) generate --template buf.gen.yaml $(PROTO_REPO)


clean:
	rm -rf services/**/internal/gen/*