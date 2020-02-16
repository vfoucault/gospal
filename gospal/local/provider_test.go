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
	"bytes"
	"context"
	"github.com/contentsquare/gospal/gospal"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
)

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
			name: "Should create a local provider",
			args: args{
				ctx:    context.Background(),
				bucket: "/tmp/bladibla",
				config: &gospal.ProviderConfig{},
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

func Test_provider_DeleteKey(t *testing.T) {

	tmpFile, err := ioutil.TempFile(os.TempDir(), "gospalTests")
	if err != nil {
		t.Errorf("unable to create temporary file for tests. err=%v", err.Error())
		return
	}
	tmpFile.WriteString("bladibla some content")
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	type fields struct {
		provider *provider
	}
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should remove the spacified key",
			fields: fields{
				provider: &provider{
					context:              context.Background(),
					kind:                 "local",
					directory:            os.TempDir(),
					noSuchKeyErrorString: os.ErrNotExist.Error(),
					config:               &gospal.ProviderConfig{},
				},
			},
			args: args{
				fileName: path.Base(tmpFile.Name()),
			},
			wantErr: false,
		},
		{
			name: "Should raise on non existing key",
			fields: fields{
				provider: &provider{
					context:              context.Background(),
					kind:                 "local",
					directory:            os.TempDir(),
					noSuchKeyErrorString: os.ErrNotExist.Error(),
					config:               &gospal.ProviderConfig{},
				},
			},
			args: args{
				fileName: "bladibla",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.provider.DeleteKey(tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("DeleteKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_provider_GetKind(t *testing.T) {
	type fields struct {
		context              context.Context
		kind                 string
		directory            string
		noSuchKeyErrorString string
		config               *gospal.ProviderConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should return the provider kind",
			fields: fields{
				context:              context.Background(),
				kind:                 "local",
				directory:            os.TempDir(),
				noSuchKeyErrorString: os.ErrNotExist.Error(),
				config:               &gospal.ProviderConfig{},
			},
			want: "local",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &provider{
				context:              tt.fields.context,
				kind:                 tt.fields.kind,
				directory:            tt.fields.directory,
				noSuchKeyErrorString: tt.fields.noSuchKeyErrorString,
				config:               tt.fields.config,
			}
			if got := p.GetKind(); got != tt.want {
				t.Errorf("GetKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_provider_GetNoSuchKeyErrorString(t *testing.T) {
	type fields struct {
		context              context.Context
		kind                 string
		directory            string
		noSuchKeyErrorString string
		config               *gospal.ProviderConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should return the proper error string",
			fields: fields{
				context:              context.Background(),
				kind:                 "local",
				directory:            os.TempDir(),
				noSuchKeyErrorString: os.ErrNotExist.Error(),
				config:               &gospal.ProviderConfig{},
			},
			want: os.ErrNotExist.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &provider{
				context:              tt.fields.context,
				kind:                 tt.fields.kind,
				directory:            tt.fields.directory,
				noSuchKeyErrorString: tt.fields.noSuchKeyErrorString,
				config:               tt.fields.config,
			}
			if got := p.GetNoSuchKeyErrorString(); got != tt.want {
				t.Errorf("GetNoSuchKeyErrorString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_provider_GetStream(t *testing.T) {

	tmpFile, err := ioutil.TempFile(os.TempDir(), "gospalTests")
	if err != nil {
		t.Errorf("unable to create temporary file for tests. err=%v", err.Error())
		return
	}
	tmpFile.WriteString("bladibla some content")
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	type fields struct {
		context              context.Context
		kind                 string
		directory            string
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
			name: "Should return the proper stream",
			fields: fields{
				context:              context.Background(),
				kind:                 "local",
				directory:            os.TempDir(),
				noSuchKeyErrorString: os.ErrNotExist.Error(),
				config:               &gospal.ProviderConfig{},
			},
			args: args{
				filePath: path.Base(tmpFile.Name()),
			},
			want:    "bladibla some content",
			wantErr: false,
		},
		{
			name: "Should raise for non exising file",
			fields: fields{
				context:              context.Background(),
				kind:                 "local",
				directory:            os.TempDir(),
				noSuchKeyErrorString: os.ErrNotExist.Error(),
				config:               &gospal.ProviderConfig{},
			},
			args: args{
				filePath: "bladibla",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &provider{
				context:              tt.fields.context,
				kind:                 tt.fields.kind,
				directory:            tt.fields.directory,
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

func Test_provider_ListKeys(t *testing.T) {

	tmpDirectory, err := ioutil.TempDir(os.TempDir(), "gospalTest")
	if err != nil {
		t.Errorf("unable to create temporary directory for tests. err=%v", err.Error())
		return
	}
	file1, err := os.Create(path.Join(tmpDirectory, "file1.txt"))
	if err != nil {
		t.Errorf("unable to create file1 for tests. err=%v", err.Error())
		return
	}
	file2, err := os.Create(path.Join(tmpDirectory, "file2.txt"))
	if err != nil {
		t.Errorf("unable to create file2 for tests. err=%v", err.Error())
		return
	}
	file3, err := os.Create(path.Join(tmpDirectory, "file3.txt"))
	if err != nil {
		t.Errorf("unable to create file3 for tests. err=%v", err.Error())
		return
	}
	err = os.MkdirAll(path.Join(tmpDirectory, "otherpath"), 0700)
	if err != nil {
		t.Errorf("unable to create otherpath directory for tests. err=%v", err.Error())
		return
	}
	file4, err := os.Create(path.Join(tmpDirectory, "otherpath/file4.txt"))
	if err != nil {
		t.Errorf("unable to create file3 for tests. err=%v", err.Error())
		return
	}
	file1.Close()
	file2.Close()
	file3.Close()
	file4.Close()
	defer os.RemoveAll(tmpDirectory)

	type fields struct {
		context              context.Context
		kind                 string
		directory            string
		noSuchKeyErrorString string
		config               *gospal.ProviderConfig
	}
	type args struct {
		pathName []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Should return all keys in local directory",
			fields: fields{
				context:              context.Background(),
				kind:                 "local",
				directory:            tmpDirectory,
				noSuchKeyErrorString: os.ErrNotExist.Error(),
				config:               &gospal.ProviderConfig{},
			},
			args: args{
				pathName: nil,
			},
			want:    []string{"/file1.txt", "/file2.txt", "/file3.txt", "/otherpath/file4.txt"},
			wantErr: false,
		},
		{
			name: "Should raise when directory does not exists",
			fields: fields{
				context:              context.Background(),
				kind:                 "local",
				directory:            "/bl/adi/bla",
				noSuchKeyErrorString: os.ErrNotExist.Error(),
				config:               &gospal.ProviderConfig{},
			},
			args: args{
				pathName: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &provider{
				context:              tt.fields.context,
				kind:                 tt.fields.kind,
				directory:            tt.fields.directory,
				noSuchKeyErrorString: tt.fields.noSuchKeyErrorString,
				config:               tt.fields.config,
			}
			got, err := p.ListKeys(tt.args.pathName...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListKeys() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_provider_PutStream(t *testing.T) {
	type fields struct {
		context              context.Context
		kind                 string
		directory            string
		noSuchKeyErrorString string
		config               *gospal.ProviderConfig
	}
	type args struct {
		fileName string
		reader   io.Reader
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        int64
		wantContent string
		wantErr     bool
	}{
		{
			name: "Should put the specified stream",
			fields: fields{
				context:              context.Background(),
				kind:                 "local",
				directory:            os.TempDir(),
				noSuchKeyErrorString: "",
				config:               &gospal.ProviderConfig{},
			},
			args: args{
				fileName: "bladibla_file1.out",
				reader:   strings.NewReader(`{"configuration": {"main_color": "#333"}, "screens": []}`),
			},
			want:        56,
			wantContent: `{"configuration": {"main_color": "#333"}, "screens": []}`,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &provider{
				context:              tt.fields.context,
				kind:                 tt.fields.kind,
				directory:            tt.fields.directory,
				noSuchKeyErrorString: tt.fields.noSuchKeyErrorString,
				config:               tt.fields.config,
			}
			got, err := p.PutStream(tt.args.fileName, tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PutStream() got = %v, want %v", got, tt.want)
			}
			if err == nil {
				// open file for comparison
				data, err := ioutil.ReadFile(path.Join(os.TempDir(), tt.args.fileName))
				if err != nil {
					t.Errorf("PutStream() unable to open local data file %v. err=%v", path.Join(os.TempDir(), tt.args.fileName), tt.want)
				}
				if string(data) != tt.wantContent {
					t.Errorf("PutStream() content = %v, wantContent %v", string(data), tt.wantContent)
				}
			}
		})
	}
}
