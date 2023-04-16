# http_response_hash

## About
A tool which makes http requests and prints the address of the request along with the MD5 hash of the response.

## Usage
```
$ cd http_response_hash
$ go build
$ ./http_response_hash [-parallel NUMBER] url1 ... urlN
```

## Notes
* There was a requirement not to include dependencies beyond Go's standard library, but `testify` package is used because
it is important to test the code properly with a proper set of tools (such as assertions for example).
* No logging was used, but in a production ready service I would use some fast and structured logging like
[Zap](https://pkg.go.dev/go.uber.org/zap).
* Parallelization limit implementation alternatives: https://gist.github.com/AntoineAugusti/80e99edfe205baf7a094
* More advanced mocks/stubs could be added in future: [golang/mock](https://github.com/golang/mock) or [testify/mock](https://pkg.go.dev/github.com/stretchr/testify/mock).