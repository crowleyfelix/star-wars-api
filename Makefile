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
	@sudo docker container stop swapi-mongo > /dev/null
	@sudo docker-compose up -d

stop: 
	@sudo docker-compose down

check: setup
	@gometalinter ./... --aggregate --fast $(MODIFIED_FILES)

deep-check: setup
	@gometalinter ./... --aggregate $(DIFF_FILES)

full-check: setup
	@gometalinter ./... --aggregate

docker-up: setup
	@sudo docker build -f docker/Dockerfile-mongo --tag star-wars-api/mongo .
	@sudo docker container start swapi-mongo || \
		sudo docker run \
		-d --rm \
		--name swapi-mongo \
		-p 27017:27017 -p 28017:28017 \
		star-wars-api/mongo
	@clear

travis:
	@gocov test -gcflags=-l --tags=integration ./... | gocov report

test: setup
	@ginkgo -gcflags=-l ./...	

test-integ: docker-up
	@ginkgo -gcflags=-l --tags=integration ./...

cov: docker-up
	@gocov test -gcflags=-l --tags=integration ./... | gocov report
	
cov-html: docker-up
	@gocov test -gcflags=-l --tags=integration ./... | gocov-html > cov.html

mock:
	@mockery -dir=./api/clients/swapi -name=Client --output=api/clients/swapi/mocks -case=underscore	
	@mockery -dir=./api/controllers -name=RequestContext --output=api/controllers/mocks -case=underscore
	@mockery -dir=./api/controllers -name=Controller --output=api/controllers/mocks -case=underscore
	@mockery -dir=./api/database/mongodb -name=SessionPool --output=api/database/mongodb/mocks -case=underscore	
	@mockery -dir=./api/database/mongodb/collections -name=Planets --output=api/database/mongodb/collections/mocks -case=underscore	
	@mockery -dir=./api/services -name=Planet --output=api/services/mocks -case=underscore	
	