.EXPORT_ALL_VARIABLES:

APP_CFG_PATH = ./configs/test/config.yaml

# APP_CFG_PATH = config.yaml

run-app:
	go run cmd/main.go 

run-app-prod:
	go run cmd/main.go --app-cfg=./configs/prod/config.yaml

run-app-test:
	go run cmd/main.go --app-cfg=./configs/test/config.yaml

build-bin:
	go build -o cmd/app ./cmd/main.go
	
build-docker-image:
	docker build -t my-app .

run-docker-container:
	docker run my-app

drop-сontainers:
	docker compose rm --stop --force