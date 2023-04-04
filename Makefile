build: go-build

go-test:
	go test -cover -coverprofile=coverage.out ./...

go-build: go-fmt go-get go-vet go-mod
	go build -o ./TrackMaster

go-fmt:
	go fmt ./...

go-vet:
	go vet ./...

go-get:
	go get ./...

go-mod:
	go mod download
