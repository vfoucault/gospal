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
	"os"
	"path"
	"syscall"
)

var (
	providerKind = ""
	bucket       = ""
	prefix       = ""
	source       = ""
	fileName     = ""
)

func main() {
	// use flags to set variable
	flag.StringVar(&providerKind, "provider", "", "the provider kind. (aws,gcp or local)")
	flag.StringVar(&bucket, "bucket", "", "the bucket name")
	flag.StringVar(&prefix, "prefix", "", "the prefix value. could be empty")
	flag.StringVar(&source, "source", "", "the source fileName.")
	flag.StringVar(&fileName, "filename", "", "the target filename. If empty we will use the source fileName")

	flag.Parse()

	// Whenever one of those three is empty, let's notice it and quit
	if providerKind == "" || bucket == "" || source == "" {
		fmt.Println("Provider and/or bucket and/or source fileName should be specified.")
		syscall.Exit(1)
	}

	// if the remote filename is empty, let's set it to the source filename
	if fileName == "" {
		fileName = path.Base(source)
	}

	// We do require a ProviderConfig
	cfg := NewProviderConfig()
	// if used a global prefix should be set
	cfg.GlobalPrefix = prefix

	// It is always nice to to have a context
	ctx := context.Background()

	var provider Gospal
	var err error

	// We use the factory here to instantiate our local provider
	if provider, err = NewProviderFactory(ctx, providerKind, bucket, cfg); err != nil {
		fmt.Printf("error creating provider %v, err=%v\n", providerKind, err.Error())
		return
	}

	// Open the local file for reading
	f, err := os.Open(source)
	defer f.Close()
	if err != nil {
		fmt.Printf("Unable to open file %v for reading. err=%v\n", source, err.Error())
		syscall.Exit(2)
	}

	// Upload it
	written, err := provider.PutStream(fileName, f)

	if err != nil {
		fmt.Println(err.Error())
		syscall.Exit(3)
	}

	fmt.Printf("local file %v streamed to %v. wrote %v bytes", source, fileName, written)
}
