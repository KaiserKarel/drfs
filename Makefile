APP_CREDS = $(shell pwd)/secrets/gdfs-308084a94a68.secret.json

.PHONY: bench lorem tests

bench:
	go test -bench=. -run=^a

lorem:
	cd examples/lorem && GOOGLE_APPLICATION_CREDENTIALS=$(APP_CREDS) go run lorem.go

tests:
	cd tests && GOOGLE_APPLICATION_CREDENTIALS=$(APP_CREDS) go test
