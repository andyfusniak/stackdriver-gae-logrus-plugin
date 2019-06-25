ODIR=./bin
VERSION=`cat VERSION`
GCLOUD_VERSION=`cat VERSION | sed 's/\./-/g'`

build:
	go build -o bin/http-service-stackdriver-logging -ldflags "-X main.version=$(VERSION)" ./examples/http-service/main.go

run:
	go run -ldflags "-X main.version=$(VERSION) -X main.projectID=$(PROJECT_ID)" ./examples/http-service/main.go

deploy-gae:
	gcloud app deploy --version=$(GCLOUD_VERSION) ./examples/http-service/main.go

clean:
	-@rm -r $(ODIR)/* 2> /dev/null || true
