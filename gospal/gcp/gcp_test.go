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
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/contentsquare/gospal/gospal"
	"github.com/fsouza/fake-gcs-server/fakestorage"
	"google.golang.org/api/iterator"
	"io"
	"reflect"
	"strings"
	"testing"
)

const testBucket = "test-bucket"

func storageInit() *storage.Client {
	server := fakestorage.NewServer([]fakestorage.Object{
		{
			BucketName: testBucket,
			Name:       "path/to/bladibla_1.txt",
			Content:    []byte("Some cool contents. with more useless chars %^&*()"),
		},
		{
			BucketName: testBucket,
			Name:       "path/two/bladibla_2.txt",
			Content:    []byte("Some cool contents. with more useless chars %^&*(), but this is not the same file"),
		},
	})
	//defer server.Stop()
	//client := server.Client()
	//return server.Client()
	//toto := server.Client().Bucket(testBucket).Objects(context.Background(), &storage.Query{
	//	//Prefix:    "",
	//	//Delimiter: "",
	//})
	//var fileList []string
	//for {
	//	attrs, err := toto.Next()
	//	if err == iterator.Done {
	//		break
	//	}
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
	//	fileList = append(fileList, attrs.Name)
	//}
	//fmt.Println(fileList)
	return server.Client()
	//object := client.Bucket("some-bucket").Object("some/object/file.txt")
	//reader, err := object.NewReader(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	//defer reader.Close()
	//data, err := ioutil.ReadAll(reader)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%s", data)
}

func TestNew(t *testing.T) {
	type args struct {
		ctx    context.Context
		bucket string
		config *gospal.ProviderConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should instantiate a GCP Client",
			args: args{
				ctx:    context.Background(),
				bucket: testBucket,
				config: &gospal.ProviderConfig{
					TimeOut: 300,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.ctx, tt.args.bucket, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(&provider{}) {
				t.Errorf("New() got = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(provider{}))
			}
		})
	}
}

func Test_provider_GetKind(t *testing.T) {

	gcpClient, err := New(context.Background(), testBucket, &gospal.ProviderConfig{
		TimeOut: 300,
	})

	if err != nil {
		t.Errorf("error when instantiating gcp client. err=%v", err.Error())
	}

	type fields struct {
		client gospal.Gospal
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should return the proper kind (gcp)",
			fields: fields{
				client: gcpClient,
			},
			want: "gcp",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.client.GetKind()
			if got != tt.want {
				t.Errorf("GetKind() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_provider_GetNoSuchKeyError(t *testing.T) {

	gcpClient, err := New(context.Background(), testBucket, &gospal.ProviderConfig{
		TimeOut: 300,
	})

	if err != nil {
		t.Errorf("error when instantiating gcp client. err=%v", err.Error())
	}

	type fields struct {
		client gospal.Gospal
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should return the proper error string",
			fields: fields{
				client: gcpClient,
			},
			want: storage.ErrObjectNotExist.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fields.client.GetNoSuchKeyErrorString()
			if got != tt.want {
				t.Errorf("GetNoSuchKeyErrorString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_provider_GetStream(t *testing.T) {
	client := storageInit()

	type fields struct {
		context              context.Context
		client               *storage.Client
		bucketName           string
		kind                 string
		noSuchKeyErrorString string
		config               *gospal.ProviderConfig
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Should upload contents on fake storage",
			fields: fields{
				context:              context.Background(),
				client:               client,
				bucketName:           testBucket,
				kind:                 "gcp",
				noSuchKeyErrorString: storage.ErrObjectNotExist.Error(),
				config: &gospal.ProviderConfig{
					TimeOut: 300,
				},
			},
			args: args{
				filePath: "path/two/bladibla_2.txt",
			},
			want:    "Some cool contents. with more useless chars %^&*(), but this is not the same file",
			wantErr: false,
		},
		{
			name: "Should raise with bad path",
			fields: fields{
				context:              context.Background(),
				client:               client,
				bucketName:           testBucket,
				kind:                 "gcp",
				noSuchKeyErrorString: storage.ErrObjectNotExist.Error(),
				config: &gospal.ProviderConfig{
					TimeOut: 300,
				},
			},
			args: args{
				filePath: "path/two/non_existing_file",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &provider{
				context:              tt.fields.context,
				client:               tt.fields.client,
				bucketName:           tt.fields.bucketName,
				kind:                 tt.fields.kind,
				noSuchKeyErrorString: tt.fields.noSuchKeyErrorString,
				config:               tt.fields.config,
			}
			got, _, err := p.GetStream(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				var bb bytes.Buffer
				io.Copy(&bb, got)
				if bb.String() != tt.want {
					t.Errorf("GetStream() got = %v, want %v", bb.String(), tt.want)
				}
			}
		})
	}
}

func Test_provider_PutStream(t *testing.T) {
	client := storageInit()

	type fields struct {
		context              context.Context
		client               *storage.Client
		bucketName           string
		kind                 string
		noSuchKeyErrorString string
		config               *gospal.ProviderConfig
	}
	type args struct {
		filePath string
		stream   io.Reader
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantWritten int64
		wantString  string
		wantErr     bool
	}{
		{
			name: "Should upload contents on fake storage",
			fields: fields{
				context:              context.Background(),
				client:               client,
				bucketName:           testBucket,
				kind:                 "gcp",
				noSuchKeyErrorString: storage.ErrObjectNotExist.Error(),
				config: &gospal.ProviderConfig{
					TimeOut: 300,
				},
			},
			args: args{
				filePath: "path/to/baldibla_3.txt",
				stream:   strings.NewReader(`bladibla, some random contents !! {cool: true}`),
			},
			wantWritten: int64(len([]byte(`bladibla, some random contents !! {cool: true}`))),
			wantString:  `bladibla, some random contents !! {cool: true}`,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &provider{
				context:              tt.fields.context,
				client:               tt.fields.client,
				bucketName:           tt.fields.bucketName,
				kind:                 tt.fields.kind,
				noSuchKeyErrorString: tt.fields.noSuchKeyErrorString,
				config:               tt.fields.config,
			}
			gotWritten, err := p.PutStream(tt.args.filePath, tt.args.stream)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWritten != tt.wantWritten {
				t.Errorf("PutStream() gotWritten = %v, want %v", gotWritten, tt.wantWritten)
			}
			// check for file presence and contents
			toto := client.Bucket(testBucket).Objects(context.Background(), &storage.Query{
				Prefix: tt.args.filePath,
			})
			var fileList []string
			for {
				attrs, err := toto.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					fmt.Println(err.Error())
				}
				fileList = append(fileList, attrs.Name)
			}
			if len(fileList) == 0 || fileList[0] != tt.args.filePath {
				t.Errorf("PutStream() file is not present on fake storage !files = %v", fileList)
			}
			// test contents
			reader, err := client.Bucket(testBucket).Object(tt.args.filePath).NewReader(context.Background())
			if err != nil {
				t.Errorf("error checking PutStream() contents err=%v", err.Error())
			} else {
				var bb bytes.Buffer
				io.Copy(&bb, reader)
				if bb.String() != tt.wantString {
					t.Errorf("error checking contents of PutStream. got = %v, want %v", bb.String(), tt.wantString)
				}

			}

		})
	}
}

func Test_provider_DeleteKey(t *testing.T) {

	client := storageInit()

	type fields struct {
		context              context.Context
		client               *storage.Client
		bucketName           string
		kind                 string
		noSuchKeyErrorString string
		config               *gospal.ProviderConfig
	}
	type args struct {
		filePath string
		stream   io.Reader
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantWritten int64
		wantString  string
		wantErr     bool
	}{
		{
			name: "Should upload contents on fake storage",
			fields: fields{
				context:              context.Background(),
				client:               client,
				bucketName:           testBucket,
				kind:                 "gcp",
				noSuchKeyErrorString: storage.ErrObjectNotExist.Error(),
				config: &gospal.ProviderConfig{
					TimeOut: 300,
				},
			},
			args: args{
				filePath: "path/two/bladibla_2.txt",
			},
			wantErr: false,
		},
		{
			name: "Should raise with unknown path",
			fields: fields{
				context:              context.Background(),
				client:               client,
				bucketName:           testBucket,
				kind:                 "gcp",
				noSuchKeyErrorString: storage.ErrObjectNotExist.Error(),
				config: &gospal.ProviderConfig{
					TimeOut: 300,
				},
			},
			args: args{
				filePath: "path/two/non_existing_file.txt",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &provider{
				context:              tt.fields.context,
				client:               tt.fields.client,
				bucketName:           tt.fields.bucketName,
				kind:                 tt.fields.kind,
				noSuchKeyErrorString: tt.fields.noSuchKeyErrorString,
				config:               tt.fields.config,
			}
			err := p.DeleteKey(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				// check for file presence
				toto := client.Bucket(testBucket).Objects(context.Background(), &storage.Query{
					Prefix: tt.args.filePath,
				})
				var fileList []string
				for {
					attrs, err := toto.Next()
					if err == iterator.Done {
						break
					}
					if err != nil {
						fmt.Println(err.Error())
					}
					fileList = append(fileList, attrs.Name)
				}
				if len(fileList) != 0 {
					for _, x := range fileList {
						if x == tt.args.filePath {
							t.Errorf("DeleteKey() file is still present on fake storage !")
						}
					}
				}
			}

		})
	}
}

func Test_provider_ListKeys(t *testing.T) {

	type fields struct {
		context              context.Context
		client               *storage.Client
		bucketName           string
		kind                 string
		noSuchKeyErrorString string
		config               *gospal.ProviderConfig
	}
	type args struct {
		pathName []string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantFileList []string
		wantErr      bool
	}{
		{
			name: "Should list keys on fake storage",
			fields: fields{
				context:              context.Background(),
				client:               storageInit(),
				bucketName:           testBucket,
				kind:                 "gcp",
				noSuchKeyErrorString: storage.ErrObjectNotExist.Error(),
				config: &gospal.ProviderConfig{
					TimeOut: 300,
				},
			},
			args: args{
				pathName: nil,
			},
			wantFileList: []string{"path/to/bladibla_1.txt", "path/two/bladibla_2.txt"},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &provider{
				context:              tt.fields.context,
				client:               tt.fields.client,
				bucketName:           tt.fields.bucketName,
				kind:                 tt.fields.kind,
				noSuchKeyErrorString: tt.fields.noSuchKeyErrorString,
				config:               tt.fields.config,
			}
			gotFileList, err := p.ListKeys(tt.args.pathName...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFileList, tt.wantFileList) {
				t.Errorf("ListKeys() gotFileList = %v, want %v", gotFileList, tt.wantFileList)
			}
		})
	}
}
