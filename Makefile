BINARY_NAME=demo_device_plugin
DOCKER_IMAGE_NAME=demo_device_plugin
DOCKER_IMAGE_LATEST="$(DOCKER_IMAGE_NAME):latest"

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) -v .

build: build-linux

img: build
	docker build -f Dockerfile -t $(DOCKER_IMAGE_LATEST) .
