VERSION   := 0.0.4-dev14
TIME      := $(shell date)
GO_MODULE := github.com/cypherfox/cloud-native-demo
LDFLAGS   := "-extldflags=-static -X '$(GO_MODULE)/pkg/version.BuildTime=$(TIME)' -X '$(GO_MODULE)/pkg/version.BuildVersion=$(VERSION)'"


docker-image: bin/bugsim
	make -C deploy/docker/bugsim VERSION=$(VERSION)

bin/bugsim: cmd/bugsim/main.go pkg/k8s/client.go cmd/bugsim/cmd/root.go cmd/bugsim/cmd/server.go
	CGO_ENABLED=0 go build -o bin/bugsim -ldflags=$(LDFLAGS) ./cmd/bugsim

helm-lint:
	docker run -it --rm  \
	    --volume $(PWD)/test/data/helm/ct.yaml:/etc/ct/ct.yaml \
		--volume $(PWD):/data \
		quay.io/helmpack/chart-testing:v3.5.0 \
	  ct lint --charts /data/deploy/helm/cloud-native-demo \
	    --chart-repos grafana=https://grafana.github.io/helm-charts \
		--chart-repos linkerd=https://helm.linkerd.io/stable \
		--config /etc/ct/ct.yaml \
		--debug --print-config

helm-package:
	make -C deploy/helm package

helm-deploy:
	step certificate create root.linkerd.cluster.local ca.crt ca.key --profile root-ca \
       --no-password --insecure
	step certificate create identity.linkerd.cluster.local issuer.crt issuer.key \
       --profile intermediate-ca --not-after 8760h --no-password --insecure \
       --ca ca.crt --ca-key ca.key
	helm upgrade cloud-native-demo ./deploy/helm/cloud-native-demo \
      --install --namespace cloud-native-demo --create-namespace --devel \
      --set-file linkerd2.identityTrustAnchorsPEM=ca.crt \
      --set-file linkerd2.identity.issuer.tls.crtPEM=issuer.crt \
      --set-file linkerd2.identity.issuer.tls.keyPEM=issuer.key \
	  --set emojivoto.namespace=cloud-native-demo

kind-load: docker-image
	docker tag github.com/cypherfox/cloud-native-demo/bugsim:$(VERSION) localhost:5001/bugsim:$(VERSION)
	kind load docker-image localhost:5001/bugsim:$(VERSION)

run-local:
	kubectl run bugsim --image=localhost:5001/bugsim:$(VERSION) --serviceaccount=bugsim-sa --expose --port=80 -- /bugsim server -p 80

stop-local:
	kubectl delete deployment bugsim
	kubectl delete service bugsim

.PHONEY: kind-load docker-image helm-lint