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
	"github.com/contentsquare/gospal/gospal/aws"
	"github.com/contentsquare/gospal/gospal/errors"
	"github.com/contentsquare/gospal/gospal/gcp"
	localprovider "github.com/contentsquare/gospal/gospal/local"
)

//NewProviderFactory constructor for a provider
func NewProviderFactory(ctx context.Context, kind string, bucket string, config *gospal.ProviderConfig) (provider gospal.Gospal, err error) {
	switch kind {
	case string(gospal.ProviderAWS):
		if provider, err = awsprovider.New(ctx, bucket, config); err != nil {
			return nil, errors.ErrorInitProvider(kind, err.Error())
		}
		return provider, err
	case string(gospal.ProviderGCP):
		if provider, err = gcpprovider.New(ctx, bucket, config); err != nil {
			return nil, errors.ErrorInitProvider(kind, err.Error())
		}
		return provider, err
	case string(gospal.ProviderLocal):
		if provider, err = localprovider.New(ctx, bucket, config); err != nil {
			return nil, errors.ErrorInitProvider(kind, err.Error())
		}
		return provider, err
	}
	return nil, errors.ErrorUnknownProvider(kind)
}
