EXECUTABLE_NAME = "{{cookiecutter.application_name}}" 
BUILD_DIR = "_build" 

run:
# create test user, can cause redundant data --  ADD_DUMMIES=1 
	COOKIE_STORE_AUTH_KEY=COOKIE_STORE_AUTH_KEY COOKIE_STORE_ENCRYPT_KEY=123456781234567812345678 RUNTIME_ENV=LOCAL CSRF_PROTECTION_KEY=1234 go run . --log-level=TRACE --port=2999

build:
	@mkdir -p $(BUILD_DIR)
# first, we build for the host OS
	go build -o _build/$(EXECUTABLE_NAME).bin .
# then, we build for Linux OS -- deployment build
	GOARCH=amd64 GOOS=linux go build -o _build/$(EXECUTABLE_NAME).linux.amd64.bin .
	du -sh _build/*

clean:
	rm -r _build/
# default database and logs
	rm /tmp/chi.db
	rm /tmp/chi.logs.*

test: build
	COOKIE_STORE_AUTH_KEY=COOKIE_STORE_AUTH_KEY COOKIE_STORE_ENCRYPT_KEY=123456781234567812345678 RUNTIME_ENV=TESTING ADD_DUMMIES=1 go test -coverprofile=test_coverage.out ./...

test-no-build: 
	RUNTIME_ENV=TESTING ADD_DUMMIES=1 go test -coverprofile=test_coverage.out ./...

test-debug:
	RUNTIME_ENV=TESTING ADD_DUMMIES=1 dlv test {{cookiecutter.project_name}}/views -test.v

coverage: test
	go tool cover -html=test_coverage.out

generate-routes: build
	_build/$(EXECUTABLE_NAME).bin --routes=true &> ./ROUTES.md

install:
	go get .
	go install honnef.co/go/tools/cmd/staticcheck@latest

lint:
	staticcheck ./...
#golangci-lint run -D errcheck ./...  => doesn't work as well

docker:
	docker build . -t skeleton-chi:local
