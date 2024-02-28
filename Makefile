# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /usr/bin/env bash -o errexit -o pipefail -o nounset

.PHONY: all build compile
all: build compile
build:
	build/build.sh
compile:
	build/compile.sh

.PHONY: dev-ui
dev-ui:
	serve -s ./ui/build

.PHONY: docker-run docker-build compose 
docker-run:
	docker run -d -it -p 3000:3000 riotpot:latest
docker-build:
	docker build -t riotpot:latest .
compose:
	docker-compose -p riotpot -f build/docker/docker-compose.yaml up -d --build