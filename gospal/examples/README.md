# Gospal implementation examples

[simple](./simple/main.go)

Simple illustration of the library by listing keys in different remote storages.

This will list keys in the local storage in path `/tmp/tmplocalbucket`

```shell script
~ # go run -tags=example . -provider local -bucket /tmp/tmplocalbucket
```

[provider2provider](./provider2provider/main.go)

provider2provider will stream a file from provider1 to provider2

The following will copy a file from a AWS bucket to a GCP bucket

```shell script
~ # AWS_REGION=eu-west-1 go run -tags=example . -dst-provider gcp -dst-bucket <GCP_BUCKET> -filename <FILENAME> -src-bucket <SOURCE AWS Bucket> -src-provider aws
```

[upload](./upload/main.go)

Simple file upload to any provider

```shell script
AWS_REGION=eu-west-1 go run -tags=example . -bucket cs.temps -source main.go -provider aws
```

[backup](./backup/main.go)

Will upload a tar of the specified directory
```shell script
~ # go run -tags=example . -provider local -bucket /tmp/tmplocalbucket -directory /tmp/toto
```