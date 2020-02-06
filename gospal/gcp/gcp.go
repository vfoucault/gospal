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

package gcpprovider

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/contentsquare/gospal/gospal"
	"github.com/contentsquare/gospal/gospal/errors"
	"google.golang.org/api/iterator"
	"io"
	"path"
	"time"
)

type provider struct {
	context              context.Context
	client               *storage.Client
	bucketName           string
	kind                 string
	noSuchKeyErrorString string

	config *gospal.ProviderConfig
}

func (p *provider) ListKeys(pathName ...string) (fileList []string, err error) {
	if len(pathName) > 1 {
		return nil, errors.ErrorTooMuchListKeysArgs()
	}
	var extraPath string
	if len(pathName) != 0 {
		extraPath = pathName[0]
	}
	targetKey := fmt.Sprintf("%v/", path.Join(p.config.GlobalPrefix, extraPath))
	it := p.client.Bucket(p.bucketName).Objects(p.context, &storage.Query{
		Prefix:    targetKey,
		Delimiter: p.config.Delimiter,
	})

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.ErrorListKeysError(targetKey, err.Error())
		}
		fileList = append(fileList, attrs.Name)

	}
	return fileList, nil
}

func (p *provider) GetStream(filePath string) (io.Reader, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(p.context, time.Second*time.Duration(p.config.TimeOut))
	var reader *storage.Reader
	var err error

	if reader, err = p.client.Bucket(p.bucketName).Object(path.Join(p.config.GlobalPrefix, filePath)).NewReader(ctx); err != nil {
		defer cancel()
		return nil, nil, errors.ErrorGetStreamReader(path.Join(p.config.GlobalPrefix, filePath), err.Error())
	}
	return reader, cancel, err
}

func (p *provider) PutStream(filePath string, stream io.Reader) (written int64, err error) {
	targetKey := path.Join(p.config.GlobalPrefix, filePath)
	ctx, cancel := context.WithTimeout(p.context, time.Second*time.Duration(p.config.TimeOut))
	defer cancel()
	wc := p.client.Bucket(p.bucketName).Object(targetKey).NewWriter(ctx)
	defer wc.Close()
	if written, err = io.Copy(wc, stream); err != nil {
		return 0, errors.ErrorPutStreamReader(targetKey, err.Error())
	}
	return written, err
}

func (p *provider) GetKind() string {
	return p.kind
}

func (p *provider) DeleteKey(filePath string) error {
	if err := p.client.Bucket(p.bucketName).Object(path.Join(p.config.GlobalPrefix, filePath)).Delete(p.context); err != nil {
		return errors.ErrorDeleteKey(path.Join(p.config.GlobalPrefix, filePath), err.Error())
	}
	return nil
}

func (p *provider) GetNoSuchKeyErrorString() string {
	return p.noSuchKeyErrorString
}

//New gcp provider constructor
func New(ctx context.Context, bucket string, config *gospal.ProviderConfig) (gospal.Gospal, error) {
	provider := provider{bucketName: bucket}
	provider.noSuchKeyErrorString = storage.ErrObjectNotExist.Error()
	provider.config = config
	provider.context = ctx
	client, err := storage.NewClient(provider.context)
	if err != nil {
		return &provider, err
	}
	provider.client = client
	return &provider, err
}
