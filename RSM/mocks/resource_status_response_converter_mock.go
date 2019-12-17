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

package mocks

import (
	"github.com/stretchr/testify/mock"
	"rsm/models"
)

type ResourceStatusResponseConverterMock struct {
	mock.Mock
}

func (m *ResourceStatusResponseConverterMock) Convert(packedBuf []byte) (*models.ResourceStatusResponse, error) {
	args := m.Called(packedBuf)
	return args.Get(0).(*models.ResourceStatusResponse), args.Error(1)
}

func (m *ResourceStatusResponseConverterMock) UnpackX2apPduAsString(packedBuf []byte, maxMessageBufferSize int) (string, error) {
	args := m.Called(packedBuf, maxMessageBufferSize)
	return args.Get(0).(string), args.Error(1)
}
