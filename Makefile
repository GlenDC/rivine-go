pkgs=$(shell go list ./... | grep -v vendor)

test-coverage:
	gocoverutil -coverprofile cover.out test \
		-short -race -v -tags='debug testing' -timeout=60s -covermode=atomic $(pkgs)

test-coverage-web: test-coverage
	go tool cover -html=cover.out
