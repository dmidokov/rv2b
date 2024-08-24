.ONESHELL:
include .env
export
all: build-and-run
build-and-run: generate-migration
	rm ./r
	/home/noname/go/go1.21.1/bin/go build -o ./r ./main.go
	./r
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
generate-migration:
	rm ./migrations/000001_tables.up.sql
	cp ./migrations/tables.sql ./migrations/000001_tables.up.sql
	sed -i 's/RV2_DB_NAME/$(DB_NAME)/g' ./migrations/000001_tables.up.sql
	sed -i 's/RV2_DOMAIN_NAME/$(BASE_DOMAIN)/g' ./migrations/000001_tables.up.sql
