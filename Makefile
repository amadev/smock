.PHONY: run
run: build
	docker run --rm -d --network host amadev/smock

.PHONY: build
build:
	docker build -t amadev/smock $(CURDIR)

.PHONY: push
push: build
	docker push amadev/smock
