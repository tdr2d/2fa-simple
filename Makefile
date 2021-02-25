TAG := $(shell git describe --tags --abbrev=0)
IMAGE := "quay.io/twebber/2fa-simple:${TAG}"

dev:
	go run .
test:
	cd tests && go test
docker_build:
	docker build . -t ${IMAGE}
docker_run:
	docker run --rm --name 2fa -p 3000:3000 ${IMAGE}
docker_sh:
	docker exec -it 2fa sh
docker_push:
	docker push ${IMAGE}

## Frontend
tailwind:
	cd web-2fa && npm run build
tailwind_prod:
	cd web-2fa && npm run build_prod


## Utils
graceful_stop:
	lsof -i :3000 | awk '{system("kill -2 " $$2)}'
clean:
	rm -rf data