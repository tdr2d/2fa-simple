run:
	go run .

clean:
	rm -rf data

test:
	cd tests && go test

graceful_stop:
	lsof -i :3000 | awk '{system("kill -2 " $$2)}'