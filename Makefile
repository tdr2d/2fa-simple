run:
	source sendgrid.env && go run .

clean:
	rm -rf data

test:
	cd tests && go test
