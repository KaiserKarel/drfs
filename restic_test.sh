export RESTIC_REPOSITORY="drfs:e2e2"
export E2E2_ACCOUNT_KEY="secretpassword"

go build -o restic-cli restic/restic/cmd/restic/*.go
./restic-cli init
./restic-cli backup README.md
