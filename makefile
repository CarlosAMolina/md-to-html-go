build:
	go build -o md-to-html .

filter:
	go test -run TestConvertFile

format:
	go fmt .

run:
	go run .

test:
	go test -v
