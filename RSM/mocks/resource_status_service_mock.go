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
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/mock"
	"rsm/models"
)

type ResourceStatusServiceMock struct {
	mock.Mock
}

func (m *ResourceStatusServiceMock) BuildAndSendInitiateRequest(nodeb *entities.NodebInfo, config *models.RsmGeneralConfiguration, enb1MeasurementId int64) error {
	args := m.Called(nodeb, config, enb1MeasurementId)
	return args.Error(0)
}
