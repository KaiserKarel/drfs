APP_CREDS = $(shell pwd)/secrets/gdfs-308084a94a68.secret.json

.PHONY: bench lorem tests restic

clean:
	rm bin

bench:
	go test -bench=. -run=^a

lorem:
	cd examples/lorem && GOOGLE_APPLICATION_CREDENTIALS=$(APP_CREDS) go run lorem.go

tests:
	cd e2e && GOOGLE_APPLICATION_CREDENTIALS=$(APP_CREDS) go test

restic:
	go build -o bin restic/restic/cmd/restic/*.go