//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package httpmsghandlerprovider

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"rsm/configuration"
	"rsm/rsmerrors"
	"rsm/tests/testhelper"
	"testing"
)

func setupTest(t *testing.T) *RequestHandlerProvider {
	config, err := configuration.ParseConfiguration()
	if err != nil {
		t.Errorf("#... - failed to parse configuration error: %s", err)
	}
	rnibDataService, rmrSender, log := testhelper.InitTestCase(t)
	return NewRequestHandlerProvider(log, rmrSender, config, rnibDataService)
}

func TestNewRequestHandlerProvider(t *testing.T) {
	provider := setupTest(t)

	assert.NotNil(t, provider)
}

func TestNewRequestHandlerProvider_InternalError(t *testing.T) {
	provider := setupTest(t)

	_, actual := provider.GetHandler("some request")

	expected := &rsmerrors.InternalError{}

	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
}
