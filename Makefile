run: generate
	# See go run . help serve for defaults.
	go run . serve

.PHONY: generate
generate:
	go get
	go generate ./...
