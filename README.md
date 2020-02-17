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

Check examples [here](./gospal/examples) 
in the example folders site a basic command line implementation:
