SHELL := /bin/bash

.PHONY: build api invoke clean

build:
	sam build

start: build
	sam local start-api -p 3001

clean:
	rm -rf .aws-sam
