TAG = 0.0.3

test:
	@docker-compose up -d --scale test=0
	@docker-compose up --abort-on-container-exit test

clean:
	@docker-compose rm -sfv

ARCH = x86_64_linux x86_64_darwin
CWD = /go/src/github.com/imega/sensible-breakpoints

release: builder $(ARCH)

builder:
	docker build -t imega/sensible-breakpoints:$(TAG) .

x86_64_%:
	@-mkdir -p build
	docker run --rm \
		-v $(CURDIR):$(CWD) \
		-w $(CWD)/cmd \
		-e GOOS=$* \
		-e GOARCH=amd64 \
    	imega/sensible-breakpoints:$(TAG) \
    	/bin/sh -c "go build -ldflags '-X main.version=$(TAG) -s' -v -o ../build/sensible-breakpoints_$@"
