IMAGE_NAME = lutzpfannenschmidt/anrechnungsstundenberechner
VERSION := $(shell date +%Y.%m.%d-%H%M%S)

buildall:
	@echo "Building all..."
	@make clean
	@mkdir -p bin/templates
	@go build -o bin/ .
	@cp -r templates/* bin/templates
	@make buildclient

buildclient:
	@echo "Building client..."
	@cd client && npm install && npm run build
	@mkdir -p bin/client/dist
	@cp -r client/dist/* bin/client/dist

clean:
	rm -rf bin/
	rm -rf client/dist

docker_build:
	@echo "Building Docker image..."
	@docker build --build-arg VERSION=$(VERSION) -t $(IMAGE_NAME):$(VERSION) .

docker_push:
	@echo "Pushing Docker image to Docker Hub..."
	@docker push $(IMAGE_NAME):$(VERSION)

publish: docker_build docker_push

docker_clean:
	@echo "Cleaning up Docker images and containers..."
	@docker rmi $(IMAGE_NAME):$(VERSION)
	@docker container prune -f

docker_test:
	@echo "Testing Docker container..."
	@docker build --build-arg VERSION=$(VERSION)_test -t $(IMAGE_NAME):test .
	@docker run -p 1337:80 $(IMAGE_NAME):test