NAMESPACE := messenger
SERVICE := contact
BUILD_PATH := ./build
K8S_PATH := ./deployments/k8s

.PHONY: build apply restart logs

build:
	eval $$(minikube docker-env); \
	docker build -t $(SERVICE):latest -f $(BUILD_PATH)/Dockerfile .

apply:
	kubectl apply -f $(K8S_PATH)/namespace.yaml
	kubectl apply -f $(K8S_PATH)/deployment_with_probes.yaml
	kubectl apply -f $(K8S_PATH)/service_cluster_ip.yaml
	@if [ -f $(K8S_PATH)/ingress.yaml ]; then \
		kubectl apply -f $(K8S_PATH)/ingress.yaml; \
	fi

restart:
	kubectl rollout restart deployment $(SERVICE)-service -n $(NAMESPACE) || true

logs:
	kubectl logs -n $(NAMESPACE) -l app=$(SERVICE)-service
