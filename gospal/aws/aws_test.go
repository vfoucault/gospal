package awsprovider

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/contentsquare/gospal/gospal"
	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3mem"
	"io"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

const testBucket = "test-bucket"

var (
	fakeS3Backend = s3mem.New()
	fakerS3       = gofakes3.New(fakeS3Backend)
	tsS3          = httptest.NewServer(fakerS3.Server())
)

func StorageReset() {
	fakeS3Backend.CreateBucket(testBucket)
	_ = os.Setenv("AWS_REGION", "eu-west-1")
	_ = os.Setenv("AWS_ENDPOINT", tsS3.URL)
}

func CreateStorageFiles() {

	type fakeDataStruct struct {
		fileName, contents string
	}
	fakeData := make([]fakeDataStruct, 0)
	fakeData = append(fakeData, fakeDataStruct{"bladibla_1.txt", `{"configuration": {"main_color": "#123"}, "screens": []}`})
	fakeData = append(fakeData, fakeDataStruct{"bladibla_2.txt", `{"configuration": {"main_color": "#345"}, "screens": [1,2]}`})
	fakeData = append(fakeData, fakeDataStruct{"bladibla/bladibla_3.out", `{"configuration": {"main_color": "#567"}, "screens": [89]}`})

	for _, data := range fakeData {
		if _, err := fakeS3Backend.PutObject(testBucket, data.fileName, make(map[string]string, 0), strings.NewReader(data.contents), int64(len([]byte(data.contents)))); err != nil {
			fmt.Println(err.Error())
		}
	}
	toto, err := fakeS3Backend.ListBucket(testBucket, &gofakes3.Prefix{}, gofakes3.ListBucketPage{})
	if err != nil {
		fmt.Println(err.Error())

	}
	for _, x := range toto.Contents {
		fmt.Println(x.Key)
	}
}

func TestNew(t *testing.T) {

	StorageReset()

	type args struct {
		ctx    context.Context
		bucket string
		config *gospal.ProviderConfig
	}
	tests := []struct {
		name     string
		extraEnv map[string]string
		args     args
		wantErr  bool
	}{
		{
			name: "Should instantiate a new aws provider",
			args: args{
				ctx:    context.Background(),
				bucket: testBucket,
				config: &gospal.ProviderConfig{
					TimeOut: 300,
					//GlobalPrefix: "/",
					//Delimiter:    "",
					MaxKeys: 1024,
				},
			},
			wantErr: false,
		},
		{
			name:     "Should raise on missing AWS_REGION",
			extraEnv: map[string]string{"AWS_REGION": "LOCAL_REGION"},
			args: args{
				ctx:    context.Background(),
				bucket: testBucket,
				config: &gospal.ProviderConfig{
					TimeOut: 300,
					//GlobalPrefix: "/",
					//Delimiter:    "",
					MaxKeys: 1024,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// process extra env
			if len(tt.extraEnv) > 0 {
				for k, v := range tt.extraEnv {
					_ = os.Setenv(k, v)
				}
			}
			got, err := New(tt.args.ctx, tt.args.bucket, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(&provider{}) {
				t.Errorf("New() got = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(provider{}))
			}
			// cleanup env
			if len(tt.extraEnv) > 0 {
				for k := range tt.extraEnv {
					_ = os.Unsetenv(k)
				}
			}
		})
	}
}

func Test_provider_ListKeys(t *testing.T) {

	StorageReset()
	CreateStorageFiles()

	awsClient, err := New(context.Background(), testBucket, &gospal.ProviderConfig{
		SpecConfig: &aws.Config{
			S3ForcePathStyle: aws.Bool(true),
		},
	})

	if err != nil {
		t.Errorf("error then setup fake s3 client. err=%v", err.Error())
	}

	delimitedAwsClient, err := New(context.Background(), testBucket, &gospal.ProviderConfig{
		SpecConfig: &aws.Config{
			S3ForcePathStyle: aws.Bool(true),
		},
		Delimiter: "/",
	})

	if err != nil {
		t.Errorf("error then setup fake s3 delimited client. err=%v", err.Error())
	}

	type fields struct {
		client gospal.Gospal
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
			name: "Should return the proper file list",
			fields: fields{
				client: awsClient,
			},
			args: args{
				pathName: []string{},
			},
			wantFileList: []string{"bladibla_1.txt", "bladibla_2.txt", "bladibla/bladibla_3.out"},
			wantErr:      false,
		},
		{
			name: "Should return the proper file list with the specified prefix",
			fields: fields{
				client: awsClient,
			},
			args: args{
				pathName: []string{"bladibla/"},
			},
			wantFileList: []string{"bladibla/bladibla_3.out"},
			wantErr:      false,
		},
		{
			name: "Should deal with delimiter",
			fields: fields{
				client: delimitedAwsClient,
			},
			args: args{
				pathName: []string{},
			},
			wantFileList: []string{"bladibla_1.txt", "bladibla_2.txt"},
			wantErr:      false,
		},
		{
			name: "Should return an error with too many paths",
			fields: fields{
				client: awsClient,
			},
			args: args{
				pathName: []string{"/", ""},
			},
			wantFileList: nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFileList, err := tt.fields.client.ListKeys(tt.args.pathName...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sort.Strings(gotFileList)
			sort.Strings(tt.wantFileList)
			if !reflect.DeepEqual(gotFileList, tt.wantFileList) {
				t.Errorf("ListKeys() gotFileList = %v, want %v", gotFileList, tt.wantFileList)
			}
		})
	}
}

func Test_provider_PutStream(t *testing.T) {

	StorageReset()

	awsClient, err := New(context.Background(), testBucket, &gospal.ProviderConfig{
		SpecConfig: &aws.Config{
			S3ForcePathStyle: aws.Bool(true),
		},
	})
	if err != nil {
		t.Errorf("error when instantiating aws client. err=%v", err.Error())
	}

	type fields struct {
		client gospal.Gospal
	}
	type args struct {
		reader   io.Reader
		fileName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "Should upload the stream",
			fields: fields{
				client: awsClient,
			},
			args: args{
				reader:   strings.NewReader(`{"configuration": {"main_color": "#333"}, "screens": []}`),
				fileName: "/bladibla_input.in",
			},
			// AWS Provider returns -1 as the number of bytes written
			want:    -1,
			wantErr: false,
		},
		{
			name: "Should reaise an error when putting the stream",
			fields: fields{
				client: awsClient,
			},
			args: args{
				reader:   strings.NewReader(`{"configuration": {"main_color": "#333"}, "screens": []}`),
				fileName: "",
			},
			// AWS Provider returns -1 as the number of bytes written
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.client.PutStream(tt.args.fileName, tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutStream() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_provider_GetStream(t *testing.T) {

	StorageReset()
	CreateStorageFiles()
	awsClient, err := New(context.Background(), testBucket, &gospal.ProviderConfig{
		TimeOut: 300,
		SpecConfig: &aws.Config{
			S3ForcePathStyle: aws.Bool(true),
		},
	})

	if err != nil {
		t.Errorf("error when instantiating aws client. err=%v", err.Error())
	}

	type fields struct {
		client gospal.Gospal
	}
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Should download the stream",
			fields: fields{
				client: awsClient,
			},
			args: args{
				fileName: "/bladibla_1.txt",
			},
			want:    `{"configuration": {"main_color": "#123"}, "screens": []}`,
			wantErr: false,
		},
		{
			name: "Should raise with unknown key",
			fields: fields{
				client: awsClient,
			},
			args: args{
				fileName: "/bladibla_1_no_such_file.txt",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := tt.fields.client.GetStream(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				var bb bytes.Buffer
				io.Copy(&bb, got)
				if bb.String() != tt.want {
					t.Errorf("GetStream() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_provider_DeleteKey(t *testing.T) {

	StorageReset()
	CreateStorageFiles()

	awsClient, err := New(context.Background(), testBucket, &gospal.ProviderConfig{
		TimeOut: 300,
		SpecConfig: &aws.Config{
			S3ForcePathStyle: aws.Bool(true),
		},
	})

	if err != nil {
		t.Errorf("error when instantiating aws client. err=%v", err.Error())
	}

	type fields struct {
		client gospal.Gospal
	}
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Should remove the specified Key",
			fields: fields{
				client: awsClient,
			},
			args: args{
				fileName: "bladibla_1.txt",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.client.DeleteKey(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// List Buket for file to check if it exists anymore
			result, err := fakeS3Backend.ListBucket(testBucket, &gofakes3.Prefix{}, gofakes3.ListBucketPage{})
			if err != nil {
				fmt.Println(fmt.Printf("error listing bucket when testing DeleteKey. err=%v", err.Error()))
			}
			for _, x := range result.Contents {
				//if x.Key == tt.args.fileName {
				if x.Key == "bladibla_1.txt" {
					t.Errorf("DeleteKey() key %v not removed.", tt.args.fileName)
				}
			}
		})
	}
}

func Test_provider_GetKind(t *testing.T) {

	StorageReset()
	awsClient, err := New(context.Background(), testBucket, &gospal.ProviderConfig{
		TimeOut: 300,
		SpecConfig: &aws.Config{
			S3ForcePathStyle: aws.Bool(true),
		},
	})

	if err != nil {
		t.Errorf("error when instantiating aws client. err=%v", err.Error())
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
			name: "Should return the proper kind (aws)",
			fields: fields{
				client: awsClient,
			},
			want: "aws",
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

	StorageReset()
	awsClient, err := New(context.Background(), testBucket, &gospal.ProviderConfig{
		TimeOut: 300,
		SpecConfig: &aws.Config{
			S3ForcePathStyle: aws.Bool(true),
		},
	})

	if err != nil {
		t.Errorf("error when instantiating aws client. err=%v", err.Error())
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
				client: awsClient,
			},
			want: s3.ErrCodeNoSuchKey,
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

func Test_provider_getTargetKey(t *testing.T) {
	type fields struct {
		context              context.Context
		session              *session.Session
		s3Service            *s3.S3
		downloader           *s3manager.Downloader
		uploader             *s3manager.Uploader
		delimiter            string
		bucketName           string
		kind                 string
		noSuchKeyErrorString string
		config               *gospal.ProviderConfig
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Should return the proper targetKey",
			fields: fields{
				config: &gospal.ProviderConfig{
					GlobalPrefix: "/bladibla",
				},
			},
			args: args{
				filePath: "test/file1.txt",
			},
			want: "/bladibla/test/file1.txt",
		},
		{
			name: "Should get ride of double slashes between paths",
			fields: fields{
				config: &gospal.ProviderConfig{
					GlobalPrefix: "/bladibla/",
				},
			},
			args: args{
				filePath: "test/file1.txt",
			},
			want: "/bladibla/test/file1.txt",
		},
		{
			name: "Should not start with slash",
			fields: fields{
				config: &gospal.ProviderConfig{
					GlobalPrefix: "",
				},
			},
			args: args{
				filePath: "test/file1.txt",
			},
			want: "test/file1.txt",
		},
		{
			name: "Should return the prefix with no trailing slash if no filename is empty",
			fields: fields{
				config: &gospal.ProviderConfig{
					GlobalPrefix: "bladitruc",
				},
			},
			args: args{
				filePath: "",
			},
			want: "bladitruc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &provider{
				context:              tt.fields.context,
				session:              tt.fields.session,
				s3Service:            tt.fields.s3Service,
				downloader:           tt.fields.downloader,
				uploader:             tt.fields.uploader,
				delimiter:            tt.fields.delimiter,
				bucketName:           tt.fields.bucketName,
				kind:                 tt.fields.kind,
				noSuchKeyErrorString: tt.fields.noSuchKeyErrorString,
				config:               tt.fields.config,
			}
			if got := p.getTargetKey(tt.args.filePath); got != tt.want {
				t.Errorf("getTargetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
