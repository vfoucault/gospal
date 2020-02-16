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

package factory

import (
	"context"
	"github.com/contentsquare/gospal/gospal"
	"os"
	"testing"
)

func TestNewProviderFactory(t *testing.T) {

	type args struct {
		ctx    context.Context
		kind   string
		bucket string
		config *gospal.ProviderConfig
	}
	tests := []struct {
		name         string
		args         args
		extraEnv     map[string]string
		wantProvider string
		wantErr      bool
	}{
		{
			name:     "Should return a aws provider",
			extraEnv: map[string]string{"AWS_REGION": "LOCAL_REGION"},
			args: args{
				ctx:    context.Background(),
				kind:   "aws",
				bucket: "test-bucket",
				config: &gospal.ProviderConfig{},
			},
			wantProvider: "aws",
			wantErr:      false,
		},
		{
			name: "Should return a gcp provider",
			args: args{
				ctx:    context.Background(),
				kind:   "gcp",
				bucket: "test-bucket",
				config: &gospal.ProviderConfig{},
			},
			wantProvider: "gcp",
			wantErr:      false,
		},
		{
			name: "Should raise on unknown provider",
			args: args{
				ctx:    context.Background(),
				kind:   "bladibla",
				bucket: "test-bucket",
				config: &gospal.ProviderConfig{},
			},
			wantProvider: "",
			wantErr:      true,
		},
		{
			name: "Should raise on missing AWS_REGION",
			args: args{
				ctx:    context.Background(),
				kind:   "bladibla",
				bucket: "test-bucket",
				config: &gospal.ProviderConfig{},
			},
			wantProvider: "aws",
			wantErr:      true,
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
			gotProvider, err := NewProviderFactory(tt.args.ctx, tt.args.kind, tt.args.bucket, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProviderFactory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if gotProvider.GetKind() != tt.wantProvider {
					t.Errorf("NewProviderFactory() gotProvider = %v, want %v", gotProvider.GetKind(), tt.wantProvider)
				}
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
