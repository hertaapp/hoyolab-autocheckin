BUILDDIR = ./build
REPO = registry-v9qyr7rcvmzk22ds.idv2.com
IMAGE_TAG = $(REPO)/hoyolab-autocheckin:latest

all: docker push_docker


docker:
	docker build -t ${IMAGE_TAG} .

push_docker:
	docker push ${IMAGE_TAG}
