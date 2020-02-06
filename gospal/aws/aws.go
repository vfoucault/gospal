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

package awsprovider

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/contentsquare/gospal/gospal"
	"github.com/contentsquare/gospal/gospal/errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type provider struct {
	context    context.Context
	session    *session.Session
	s3Service  *s3.S3
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
	delimiter  string

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
	targetKey := fmt.Sprintf("%v", path.Join(p.config.GlobalPrefix, extraPath))
	fmt.Println("targetkey is ", targetKey)
	params := s3.ListObjectsInput{
		Bucket:    aws.String(p.bucketName),
		Prefix:    aws.String(targetKey),
		Delimiter: aws.String(p.config.Delimiter),
		MaxKeys:   aws.Int64(p.config.MaxKeys),
	}

	page := request.Pagination{
		NewRequest: func() (*request.Request, error) {
			req, _ := p.s3Service.ListObjectsRequest(&params)
			req.SetContext(p.context)
			return req, nil
		},
	}

	for page.Next() {
		page := page.Page().(*s3.ListObjectsOutput)
		for _, obj := range page.Contents {
			if !strings.HasSuffix(*obj.Key, "/") {
				fileList = append(fileList, *obj.Key)
			}
		}
	}
	return
}

func (p *provider) GetStream(filePath string) (io.Reader, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(p.context, time.Second*time.Duration(p.config.TimeOut))
	targetKey := path.Join(p.config.GlobalPrefix, filePath)
	result, err := p.s3Service.GetObjectWithContext(ctx,
		&s3.GetObjectInput{
			Bucket: &p.bucketName,
			Key:    &targetKey,
		})
	if err != nil {
		defer cancel()
		return nil, nil, errors.ErrorGetStreamReader(targetKey, err.Error())
	}
	return result.Body, cancel, nil
}

func (p *provider) PutStream(filePath string, reader io.Reader) (int64, error) {
	targetKey := path.Join(p.config.GlobalPrefix, filePath)
	_, err := p.uploader.Upload(&s3manager.UploadInput{
		Bucket: &p.bucketName,
		Key:    &targetKey,
		Body:   reader,
	})
	if err != nil {
		return 0, errors.ErrorPutStreamReader(filepath.Join(p.config.GlobalPrefix, filePath), err.Error())
	}
	return -1, nil
}

func (p *provider) GetKind() string {
	return p.kind
}

func (p *provider) DeleteKey(filePath string) error {
	targetKey := path.Join(p.config.GlobalPrefix, filePath)
	if _, err := p.s3Service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &p.bucketName,
		Key:    &targetKey,
	}); err != nil {
		return errors.ErrorDeleteKey(path.Join(p.config.GlobalPrefix, filePath), err.Error())
	}
	if err := p.s3Service.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: &p.bucketName,
		Key:    &targetKey,
	}); err != nil {
		return errors.ErrorDeleteKey(path.Join(p.config.GlobalPrefix, filePath), err.Error())
	}
	return nil
}

func (p *provider) GetNoSuchKeyErrorString() string {
	return p.noSuchKeyErrorString
}

//New aws provider constructor
func New(ctx context.Context, bucket string, config *gospal.ProviderConfig) (gospal.Gospal, error) {
	// fetch aws region from env
	var region string

	if region = os.Getenv("AWS_REGION"); region == "" {
		return nil, errors.ErrorInitProvider("aws", "AWS_REGION is not set")
	}
	provider := provider{bucketName: bucket}
	provider.noSuchKeyErrorString = s3.ErrCodeNoSuchKey
	provider.config = config
	provider.context = ctx
	provider.kind = "aws"

	cfg := &aws.Config{
		// TODO: make this a confugration setting settable a test time !
		S3ForcePathStyle: aws.Bool(true),
	}

	if endpoint := os.Getenv("AWS_ENDPOINT"); endpoint != "" {
		cfg.Endpoint = aws.String(endpoint)
	}

	provider.session = session.Must(session.NewSession(cfg))
	provider.s3Service = s3.New(provider.session)
	provider.uploader = s3manager.NewUploader(provider.session)
	provider.downloader = s3manager.NewDownloader(provider.session)

	return &provider, nil
}
