IMG_NAME=project-spiderweb-service
TAG=local
VOLUME_TO_MOUNT=$(PWD)
INTERNAL_PORT=8782
EXTERNAL_PORT=$(INTERNAL_PORT)

IMG=sharded.cards/$(IMG_NAME)
CONTAINER=$(IMG_NAME)
VOLUME_DESTINATION=/go/src/github.com/Pergamene/project-spiderweb-service
STATIC_PATH=static

build:
	docker build --pull \
	-t $(IMG):$(TAG) \
	-f Dockerfile .

run-local: rm
	docker run -d \
	-p $(EXTERNAL_PORT):$(INTERNAL_PORT) \
	--name $(CONTAINER) \
	-v $(VOLUME_TO_MOUNT):$(VOLUME_DESTINATION) \
	-e STATIC_PATH=$(STATIC_PATH) \
	-e DATACENTER=LOCAL \
	-e ENVIRONMENT=local \
	-e PORT=$(INTERNAL_PORT) \
	$(IMG):$(TAG)

rm:
	docker rm \
	-f $(CONTAINER) || true

test:
	docker run \
	--name $(CONTAINER)_test --rm  \
	-e STATIC_PATH=$(STATIC_PATH) \
	-e DATACENTER=LOCAL \
	-e ENVIRONMENT=local \
	$(IMG):$(TAG) \
	go test ./...

# tags the image with the most recent git commit
tag:
	$(eval TAG := $(shell git rev-parse --verify HEAD))