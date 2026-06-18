VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/DevShedLabs/new/cmd.Version=$(VERSION)"

.PHONY: build install clean

build:
	go build $(LDFLAGS) -o new .

install:
	go install $(LDFLAGS) .

clean:
	rm -f new
