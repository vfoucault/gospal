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

package gospal

import (
	"context"
	"io"
)

const (
	defaultTimeOut = 300
	defaultMaxKeys = 1024
)

//ProviderLabel represents a provider string label such as gcp or aws
type ProviderLabel string

var (
	//ProviderAWS ProviderLabel for AWS
	ProviderAWS ProviderLabel = "aws"
	//ProviderGCP ProviderLabel for GCP
	ProviderGCP ProviderLabel = "gcp"
)

//Gospal interface that represents a Storage Gospal
//  * ListKeys: List all keys in a specified optional path within the configured bucket
//  * GetStream: Stream out a specified. Return a io.Reader of the specified key within the configured bucket
//  * PutStream: Stream in a given io.Reader to the specified key within the configured bucket
//  * GetKind: Return the provider kind, the provider name
//  * DeleteKey: remove the specified key within the configured bucket
//  * GetNoSuchKeyErrorString: return the error message for this provider when a key is not found
type Gospal interface {
	ListKeys(...string) ([]string, error)
	GetStream(string) (io.Reader, context.CancelFunc, error)
	PutStream(string, io.Reader) (int64, error)

	GetKind() string

	DeleteKey(string) error

	GetNoSuchKeyErrorString() string
}

// ProviderConfig holds common configuration between providers
type ProviderConfig struct {
	// Timeout value for the context.timeout for some operations
	TimeOut int

	// A global prefix to be applied to all path queries. eg. when not set
	// for aws will call s3://bucketName/keyName
	// when set to bladibla
	// for aws will call s3://bucketName/bladibla/keyName
	// A single character used to separate individual fields in a record. You can
	GlobalPrefix string

	// specify an arbitrary delimiter.
	Delimiter string

	// AWS Only: max number of key to fetch at once
	MaxKeys int64

	// Provider Specific Configuration.
	// type of interface to mach any, but will be reflected to specific provider configuration.
	// eg.: &aws.Config
	SpecConfig interface{}
}

//NewProviderConfig constructor with default value setter
func NewProviderConfig() *ProviderConfig {
	// TimeOut is set to 300 seconds by default
	// MaxKeys is set to 1024 by default
	return &ProviderConfig{TimeOut: defaultTimeOut, MaxKeys: defaultMaxKeys}
}
