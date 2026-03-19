build:
	go build -o md-to-html .

filter:
	go test -run TestConvertLines

format:
	go fmt .

run:
	cp -r testdata/dir-to-convert /tmp/ && go run . /tmp/dir-to-convert

test:
	go test -v
