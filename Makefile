ifndef NOTGCP
  PROJECT_ID := $(shell gcloud config get-value project)
  ZONE := $(shell gcloud config get-value compute/zone)
  SHORT_SHA := $(shell git rev-parse --short HEAD)
  IMG ?= gcr.io/${PROJECT_ID}/airflow-operator:${SHORT_SHA}
endif


# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

all: test manager

# Run tests
test: generate fmt vet manifests
	go test ./controllers/... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
debug: generate fmt vet
	dlv debug cmd/manager/main.go

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | kubectl apply -f -
	kubectl apply -f hack/appcrd.yaml

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | kubectl apply -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
undeploy: manifests
	kustomize build config/default | kubectl delete -f -
	kubectl delete -f config/crds || true
	kubectl delete -f hack/appcrd.yaml || true

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths="./..."

# Build the docker image
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	cd config/manager && kustomize edit set image controller=${IMG}

# Push the docker image
docker-push: docker-build
	docker push ${IMG}


e2e-test:
	kubectl get namespace airflowop-system || kubectl create namespace airflowop-system
	go test -v -timeout 20m test/e2e/base/base_test.go --namespace airflowop-system
	go test -v -timeout 20m test/e2e/cluster/cluster_test.go --namespace airflowop-system

e2e-test-gcp:
	kubectl get namespace airflowop-system || kubectl create namespace airflowop-system
	kubectl apply -f hack/sample/cloudsql-celery/sqlproxy-secret.yaml -n airflowop-system
	go test -v -timeout 20m test/e2e/gcp/gcp_test.go --namespace airflowop-system


# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.4 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
