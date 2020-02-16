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

package localprovider

import (
	"context"
	"fmt"
	"github.com/contentsquare/gospal/gospal"
	"github.com/contentsquare/gospal/gospal/errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type provider struct {
	context              context.Context
	kind                 string
	directory            string
	noSuchKeyErrorString string
	config               *gospal.ProviderConfig
}

func (p *provider) ListKeys(pathName ...string) ([]string, error) {
	if len(pathName) > 1 {
		return nil, errors.ErrorTooMuchListKeysArgs()
	}
	var extraPath string
	if len(pathName) != 0 {
		extraPath = pathName[0]
	}
	var files []string
	err := filepath.Walk(path.Join(p.directory, extraPath), func(filePath string, info os.FileInfo, err error) error {
		if err == nil {
			if !info.IsDir() {
				files = append(files, strings.Replace(filePath, p.directory, "", 1))
			}
		} else {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (p *provider) PutStream(fileName string, reader io.Reader) (int64, error) {
	// create a file with the proper mode. Whenever a file exists with the same name we will overwrite it
	fh, err := os.OpenFile(path.Join(p.directory, fileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer fh.Close()
	if err != nil {
		return -1, fmt.Errorf("unable to open file %v for writing. err=%v", path.Join(p.directory, fileName), err.Error())
	}
	written, err := io.Copy(fh, reader)
	if err != nil {
		return -1, fmt.Errorf("unable to write file %v. err=%v", path.Join(p.directory, fileName), err.Error())
	}
	return written, err
}

func (p *provider) GetKind() string {
	return p.kind
}

func (p *provider) DeleteKey(fileName string) error {
	// check if key exists
	_, err := os.Stat(path.Join(p.directory, fileName))
	if err != nil {
		return errors.ErrorDeleteKey(path.Join(p.directory, fileName), err.Error())
	}
	err = os.Remove(path.Join(p.directory, fileName))
	if err != nil {
		return errors.ErrorDeleteKey(path.Join(p.directory, fileName), err.Error())
	}
	return nil
}

func (p *provider) GetNoSuchKeyErrorString() string {
	return p.noSuchKeyErrorString
}

func (p *provider) GetStream(filePath string) (io.Reader, context.CancelFunc, error) {
	// the context will be useless here, let's just create a cancel func for match the method signature in the interface
	_, cancel := context.WithCancel(p.context)
	// fetch the specified file from the local filesystem
	fh, err := os.Open(path.Join(p.directory, filePath))

	// if the file does not exists, then and error should be raised
	if err != nil {
		defer cancel()
		return nil, nil, fmt.Errorf("could not open file %v. err=%v", filePath, err.Error())
	}

	// we may safely return a io.Reader from the file handler. *File implements the interface io.Reader
	return fh, cancel, nil
}

//New aws provider constructor
func New(ctx context.Context, bucket string, config *gospal.ProviderConfig) (gospal.Gospal, error) {
	provider := provider{directory: bucket}
	// check to directory to exists
	if _, err := os.Stat(bucket); os.IsNotExist(err) {
		if err2 := os.MkdirAll(bucket, 0600); err2 != nil {
			return nil, fmt.Errorf("unable to create local directory %v. err=%v", bucket, err.Error())
		}
	}
	provider.noSuchKeyErrorString = os.ErrNotExist.Error()
	provider.config = config
	provider.context = ctx
	provider.kind = "local"

	return &provider, nil
}
