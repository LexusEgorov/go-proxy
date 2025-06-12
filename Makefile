IMAGE_NAME=go-proxy

local:
	go run ./cmd/proxy --config ./configs/local.yml
build:
	docker build -t $(IMAGE_NAME) .
run:
	docker run -p 8081:8081 $(IMAGE_NAME) --config ./configs/local.yml