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

package errors

import (
	"fmt"
)

const (
	tooMuchListKeysArgsMessage        = "ListKeys takes at most one path. extra=%v"
	listKeysErrorMessage              = "ListKeys error when listing remote storage %v. extra=%v"
	getStreamReaderErrorMessage       = "GetStream: error while fetching reader for filePath %v. err=%v"
	putStreamReaderErrorMessage       = "PutStream: error when putting stream to file %v. err=%v"
	deleteKeyErrorMessage             = "DeleteKey: error when deleting key %v. err=%v"
	providerFactoryInitErrorMessage   = "NewProviderFactory: error when instantiating provider %v. err=%v"
	providerFactoryUnknownKindMessage = "NewProviderFactory: unable to process ConfigFactory. Unknown provider %v"
)

//ErrorTooMuchListKeysArgs helper to return a common error message when a too much args are given for the list function
func ErrorTooMuchListKeysArgs(extra ...interface{}) error {
	return fmt.Errorf(tooMuchListKeysArgsMessage, extra...)
}

//ErrorListKeysError helper to return a common error message when error occures in ListKeys
func ErrorListKeysError(extra ...interface{}) error {
	return fmt.Errorf(listKeysErrorMessage, extra...)
}

//ErrorGetStreamReader helper to return a common error message when an error is raised when fetching a stream from object storage
func ErrorGetStreamReader(extra ...interface{}) error {
	return fmt.Errorf(getStreamReaderErrorMessage, extra...)
}

//ErrorPutStreamReader helper to return a common error message when an error is raised when putting a stream to object storage
func ErrorPutStreamReader(extra ...interface{}) error {
	return fmt.Errorf(putStreamReaderErrorMessage, extra...)
}

//ErrorDeleteKey helper to return a common error message when an error is raised when removing a key from the object storage
func ErrorDeleteKey(extra ...interface{}) error {
	return fmt.Errorf(deleteKeyErrorMessage, extra...)
}

//ErrorInitProvider helper to return a common error message when an error is raised when instantiating a provider client
func ErrorInitProvider(extra ...interface{}) error {
	return fmt.Errorf(providerFactoryInitErrorMessage, extra...)
}

//ErrorUnknownProvider helper to return a common error message when an error is raised when instantiating a provider client for an unknown provider type
func ErrorUnknownProvider(extra ...interface{}) error {
	return fmt.Errorf(providerFactoryUnknownKindMessage, extra...)
}
