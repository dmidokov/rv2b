.ONESHELL:
include .env
export
all: build-and-run
build-and-run:
	/home/noname/go/go1.21.1/bin/go build ./main.go
	./main
build-compose:
	docker-compose down
	docker rmi rv2
	cd ../front_v/rv2/
	yarn build
	cp -R ./dist ../../backend/docker/
	cd ../../backend/
	docker-compose up
build-with-local-front:
	docker-compose down
	docker rmi rv2
	cd ../front_v/rv2/
	yarn build
	cp -R ./dist ../../backend/docker/
	cd ../../backend/
	rm ./main
	docker-compose -f ./docker-compose-with-local-front.yml up
up-with-local:
	docker-compose down
	docker-compose -f ./docker-compose-with-local-front.yml up
up:
	docker-compose down
	docker-compose up
