VERSION="0.0.2"
IMAGE_NAME="bugsim"
REGISTRY="github.com/cypherfox/cloud-native-demo"

build-image:
	cp ../../../bin/bugsim .
	docker build -t $(REGISTRY)/$(IMAGE_NAME):$(VERSION) .
	docker push $(REGISTRY)/$(IMAGE_NAME):$(VERSION)
