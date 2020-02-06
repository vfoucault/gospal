//  Copyright 2019 Contentsquare
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

// +build example

package main

import (
	"context"
	"flag"
	"fmt"
	. "github.com/contentsquare/gospal/gospal"
	. "github.com/contentsquare/gospal/gospal/factory"
	"syscall"
)

var (
	srcProviderKind = ""
	srcBucket       = ""
	dstProviderKind = ""
	dstBucket       = ""
	filename        = ""
)

func main() {

	flag.StringVar(&srcProviderKind, "src-provider", "", "the source provider kind. either aws or gcp")
	flag.StringVar(&srcBucket, "src-bucket", "", "the source bucket name")
	flag.StringVar(&dstProviderKind, "dst-provider", "", "the target provider kind. either aws or gcp")
	flag.StringVar(&dstBucket, "dst-bucket", "", "the target bucket name")
	flag.StringVar(&filename, "filename", "", "the filename")

	flag.Parse()

	fmt.Println(srcProviderKind, srcBucket, dstProviderKind, dstBucket, filename)
	if srcProviderKind == "" || srcBucket == "" || filename == "" {
		fmt.Println("source Provider/bucket/file should be specified.")
		syscall.Exit(1)
	}

	if dstProviderKind == "" || dstBucket == "" {
		fmt.Println("target Provider/bucket should be specified.")
		syscall.Exit(1)
	}

	cfg := NewProviderConfig()

	ctx := context.Background()

	var srcProvider Gospal
	var dstProvider Gospal
	var err error

	if srcProvider, err = NewProviderFactory(ctx, srcProviderKind, srcBucket, cfg); err != nil {
		fmt.Println(fmt.Sprintf("error creating provider %v, err=%v", srcProviderKind, err.Error()))
		return
	}

	if dstProvider, err = NewProviderFactory(ctx, dstProviderKind, dstBucket, cfg); err != nil {
		fmt.Println(fmt.Sprintf("error creating provider %v, err=%v", dstProviderKind, err.Error()))
		return
	}

	srcStream, _, err := srcProvider.GetStream(filename)
	if err != nil {
		fmt.Println(fmt.Sprintf("error fetching source stream %v. err=%v", filename, err.Error()))
		syscall.Exit(1)
	}

	_, err = dstProvider.PutStream(filename, srcStream)
	if err != nil {
		fmt.Println(fmt.Sprintf("error puting stream %v. err=%v", filename, err.Error()))
		syscall.Exit(1)
	}

}
