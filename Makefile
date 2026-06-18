.PHONY: build install clean

build:
	go build -o new .

install:
	go install .

clean:
	rm -f new
