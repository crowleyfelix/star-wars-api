DIFF_FILES=$(shell git diff-index --cached --name-only origin/develop | xargs printf -- '--include=%s\n')
MODIFIED_FILES=$(shell git diff-index --cached --name-only HEAD | xargs printf -- '--include=%s\n')

.PHONY: dep setup test coverage mocks

dep:
	@go get -t -v ./... \
	github.com/onsi/ginkgo/ginkgo \
	github.com/onsi/gomega/...  \
	github.com/axw/gocov/gocov \
	gopkg.in/matm/v1/gocov-html \
	github.com/vektra/mockery/.../ \
	github.com/alecthomas/gometalinter
	@gometalinter --install > /dev/null

fmt:
	@go fmt ./...

setup: dep

run: dep
	@go run main.go

check: setup
	@gometalinter ./... --aggregate --fast $(MODIFIED_FILES)

deep-check: setup
	@gometalinter ./... --aggregate $(DIFF_FILES)

full-check: setup
	@gometalinter ./... --aggregate

test: setup
	@ginkgo -gcflags=-l ./...	

test-integ:
	@ginkgo -gcflags=-l --tags=integration ./...

cov: setup
	@gocov test -gcflags=-l --tags=integration ./... | gocov report
	
cov-html: setup
	@gocov test -gcflags=-l --tags=integration ./... | gocov-html > cov.html

mock:
	@mockery -dir=./src/controllers -name=RequestContext --output=src/controllers/mocks -case=underscore
	@mockery -dir=./src/controllers -name=Controller --output=src/controllers/mocks -case=underscore
	@mockery -dir=./src/mongodb -name=SessionManager --output=src/mongodb/mocks -case=underscore	
	