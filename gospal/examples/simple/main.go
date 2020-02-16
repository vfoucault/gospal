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
	"path"
	"syscall"
)

var (
	providerKind = ""
	bucket       = ""
	prefix       = ""
)

func main() {

	flag.StringVar(&providerKind, "provider", "", "the provider kind. (aws,gcp or local)")
	flag.StringVar(&bucket, "bucket", "", "the bucket name")
	flag.StringVar(&prefix, "prefix", "", "the prefix value. could be empty")

	flag.Parse()

	if providerKind == "" || bucket == "" {
		fmt.Println("Provider and/or bucket should be specified.")
		syscall.Exit(1)
	}

	cfg := NewProviderConfig()
	cfg.GlobalPrefix = prefix

	ctx := context.Background()

	var provider Gospal
	var err error

	if provider, err = NewProviderFactory(ctx, providerKind, bucket, cfg); err != nil {
		fmt.Println(fmt.Sprintf("error creating provider %v, err=%v", providerKind, err.Error()))
		return
	}
	keys, err := provider.ListKeys()

	if err != nil {
		fmt.Println(err.Error())
		syscall.Exit(1)
	}
	fmt.Println(fmt.Sprintf("Keys in path %v :", path.Join(bucket, cfg.GlobalPrefix)))
	for _, k := range keys {
		fmt.Println(fmt.Sprintf("\t* %v", k))
	}

}
