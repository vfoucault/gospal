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

package main

import (
	"archive/tar"
	"context"
	"flag"
	"fmt"
	. "github.com/contentsquare/gospal/gospal"
	. "github.com/contentsquare/gospal/gospal/factory"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	providerKind = ""
	bucket       = ""
	prefix       = ""
	directory    = ""
)

func main() {

	flag.StringVar(&providerKind, "provider", "", "the provider kind. (aws,gcp or local)")
	flag.StringVar(&bucket, "bucket", "", "the bucket name")
	flag.StringVar(&prefix, "prefix", "", "the prefix value. could be empty")
	flag.StringVar(&directory, "directory", "", "the directory to upload")

	flag.Parse()

	if providerKind == "" || bucket == "" || directory == "" {
		fmt.Println("Provider and/or bucket and/or directory should be specified.")
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

	// io.Pipe is a really conveniant way to link a reader and a writer
	reader, writer := io.Pipe()

	// let's launch a go routing that will tar the directory.
	// this go routine will return (wg.Done()) one it is done for it to add all required files
	go tarThis(directory, writer)

	//
	written, err := provider.PutStream(path.Base(directory)+".tar", reader)
	if err != nil {
		fmt.Printf("error writing %v.tar for remote storage. err=%v", path.Base(directory), err.Error())
	}
	fmt.Printf("wrote %v bytes to remote file %v.tar.", written, directory)
}

func tarThis(directory string, writer io.WriteCloser) {
	var files = make(map[string]os.FileInfo)
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			files[path] = info
		}
		return nil
	})
	if err != nil {
		fmt.Printf(err.Error() + "\n")
		return
	}

	tarfileWriter := tar.NewWriter(writer)
	for fpath, fileInfo := range files {
		if fileInfo.Mode().IsRegular() {
			file, err := os.Open(fpath)
			if err != nil {
				fmt.Printf("Error opening %v for reading. %v\n", fpath, err)
				break
			}
			// prepare the tar header

			header := new(tar.Header)
			header.Name = path.Join(path.Base(directory), strings.TrimPrefix(file.Name(), path.Clean(directory)))
			header.Size = fileInfo.Size()
			header.Mode = int64(fileInfo.Mode())
			header.ModTime = fileInfo.ModTime()

			err = tarfileWriter.WriteHeader(header)
			if err != nil {
				fmt.Printf("Error writing tar header for %v. %v\n", fpath, err)
				file.Close()
				break
			}
			_, err = io.Copy(tarfileWriter, file)
			if err != nil {
				fmt.Printf("Error copying tar stream to tar writer for %v. %v\n", fpath, err)
				file.Close()
				break
			}
			file.Close()
		}
	}
	tarfileWriter.Close()
	writer.Close()
}
