run:
	go run .

tailwind:
	cd web-2fa && npm run build

tailwind_prod:
	cd web-2fa && npm run build_prod

clean:
	rm -rf data

test:
	cd tests && go test

graceful_stop:
	lsof -i :3000 | awk '{system("kill -2 " $$2)}'