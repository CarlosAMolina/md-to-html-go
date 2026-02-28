build:
	go build -o md-to-html .

filter:
	go test -run TestConvertFile

format:
	go fmt .

run:
	cp -r testdata/md/ /tmp/ && go run . /tmp/md

test:
	go test -v
