NAMESPACE := messenger
SERVICE := auth-service
BUILD_PATH := ./build
K8S_PATH := ./deployments/k8s

PROTO_REPO=../../api/proto
BUF=buf


.PHONY: build apply restart logs proto clean

build:
	eval $$(minikube docker-env); \
	docker build -t $(SERVICE):latest -f $(BUILD_PATH)/Dockerfile .

apply:
	kubectl apply -f $(K8S_PATH)/namespace.yaml
	kubectl apply -f $(K8S_PATH)/deployment_with_probes.yaml
	kubectl apply -f $(K8S_PATH)/service_cluster_ip.yaml
	kubectl apply -f $(K8S_PATH)/ingress_grpc.yaml
	kubectl apply -f $(K8S_PATH)/ingress_http.yaml

restart:
	kubectl rollout restart deployment $(SERVICE) -n $(NAMESPACE) || true

logs:
	kubectl logs -n $(NAMESPACE) -l app=$(SERVICE)

proto:
	$(BUF) generate --template buf.gen.yaml $(PROTO_REPO)

clean:
	rm -rf ../services/**/internal/gen/*