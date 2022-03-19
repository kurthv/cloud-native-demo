VERSION   := 0.0.4-dev14
TIME      := $(shell date)
GO_MODULE := github.com/cypherfox/cloud-native-demo
LDFLAGS   := "-extldflags=-static -X '$(GO_MODULE)/pkg/version.BuildTime=$(TIME)' -X '$(GO_MODULE)/pkg/version.BuildVersion=$(VERSION)'"


docker-image: bin/bugsim
	make -C deploy/docker/bugsim VERSION=$(VERSION)

bin/bugsim: cmd/bugsim/main.go pkg/k8s/client.go cmd/bugsim/cmd/root.go cmd/bugsim/cmd/server.go
	CGO_ENABLED=0 go build -o bin/bugsim -ldflags=$(LDFLAGS) ./cmd/bugsim

helm-lint:
	docker run -it --rm -v $(PWD):/data quay.io/helmpack/chart-testing:v3.5.0 \
	  ct lint --charts /data/deploy/helm/cloud-native-demo \
	    --chart-repos grafana=https://grafana.github.io/helm-charts \
		--chart-repos linkerd=https://helm.linkerd.io/stable \
		--debug --validate-maintainers=false

helm-package:
	make -C deploy/helm package

helm-deploy:
	helm upgrade cloud-native-demo ./deploy/helm/cloud-native-demo --install --namespace cloud-native-demo --create-namespace --devel

kind-load: docker-image
	docker tag github.com/cypherfox/cloud-native-demo/bugsim:$(VERSION) localhost:5001/bugsim:$(VERSION)
	kind load docker-image localhost:5001/bugsim:$(VERSION)

run-local:
	kubectl run bugsim --image=localhost:5001/bugsim:$(VERSION) --serviceaccount=bugsim-sa --expose --port=80 -- /bugsim server -p 80

stop-local:
	kubectl delete deployment bugsim
	kubectl delete service bugsim

.PHONEY: kind-load docker-image helm-lint