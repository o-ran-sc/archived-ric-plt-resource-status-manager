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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).


package httpmsghandlerprovider

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"rsm/configuration"
	"rsm/handlers/httpmsghandlers"
	"rsm/logger"
	"rsm/mocks"
	"rsm/rsmerrors"
	"rsm/services"
	"testing"
)

func setupTest(t *testing.T) *RequestHandlerProvider {
	log, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}

	config, err := configuration.ParseConfiguration()
	if err != nil {
		t.Errorf("#... - failed to parse configuration error: %s", err)
	}

	resourceStatusServiceMock := &mocks.ResourceStatusServiceMock{}
	rnibReaderMock := &mocks.RnibReaderMock{}
	rsmReaderMock := &mocks.RsmReaderMock{}
	rsmWriterMock := &mocks.RsmWriterMock{}

	rnibDataService := services.NewRnibDataService(log, config, rnibReaderMock, rsmReaderMock, rsmWriterMock)
	return NewRequestHandlerProvider(log, rnibDataService, resourceStatusServiceMock)
}

func TestNewRequestHandlerProvider(t *testing.T) {
	provider := setupTest(t)

	assert.NotNil(t, provider)
}

func TestResourceStatusRequestHandler(t *testing.T) {
	provider := setupTest(t)
	handler, err := provider.GetHandler(ResourceStatusRequest)

	assert.NotNil(t, provider)
	assert.Nil(t, err)

	_, ok := handler.(*httpmsghandlers.ResourceStatusRequestHandler)

	assert.True(t, ok)
}

func TestNewRequestHandlerProvider_InternalError(t *testing.T) {
	provider := setupTest(t)

	_, actual := provider.GetHandler("some request")

	expected := &rsmerrors.InternalError{}

	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
}