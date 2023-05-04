BUILDDIR = ./build
IMAGE_TAG = gcr.io/aki149/hoyolab-autocheckin:latest

.PHONY: build clean

docker:
	docker build -t ${IMAGE_TAG} .

push_docker:
	docker push ${IMAGE_TAG}
