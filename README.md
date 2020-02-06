# Gospal: Go-Object Storage Pal

disclaimer: this code intend to show the implementation of the factory model in goLang.

## Useage

```
import "github.com/contentsquare/gospal"
```

## Description

A very simple Multi csp object streaming storage client with the following methods:

```go
* ListKeys(...string) ([]string, error)
* GetStream(string) (io.Reader, context.CancelFunc, error)
* PutStream(string, io.Reader) (int64, error)
* GetKind() string
* DeleteKey(string) error
* GetNoSuchKeyErrorString() string
```

# Basic usage

in the example folders site a basic command line implementation:

```shell script
cd examples/simple

AWS_REGION=eu-west-1 go run -tags=example . -provider aws -bucket <bucket-name> -prefix <prefix-value>
```

* [simple example](./gospal/examples/simple/main.go)
* [provider2provider example](./gospal/examples/provider2provider/main.go)

