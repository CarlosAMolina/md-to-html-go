build:
	go build -o md-to-html .

coverage:
	go test -coverprofile=coverage.out && go tool cover -html=coverage.out

deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

filter:
	go test -run TestConvertLines

format:
	go fmt .

lint:
	golangci-lint run

run:
	cp -r testdata/dir-to-convert /tmp/ && go run . /tmp/dir-to-convert

test:
	go test -v
