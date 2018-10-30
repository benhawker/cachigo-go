build:
	cd cmd/api/ && go build -o build/cachigo

run:
	cmd/api/build/cachigo

test:
	go test -cover -race `go list ./... | grep -v /vendor`

install:
	dep init

update:
	dep ensure

swagger-docs:
	swagger-codegen generate -i swagger.yml -l html