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

package rmrmsghandlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"rsm/configuration"
	"rsm/e2pdus"
	"rsm/logger"
	"rsm/mocks"
	"rsm/models"
	"testing"
	"time"
)

func initRanConnectedNotificationHandlerTest(t *testing.T, requestName string) (ResourceStatusInitiateNotificationHandler, *mocks.ResourceStatusInitiateManagerMock, *configuration.Configuration) {
	log, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	config, err := configuration.ParseConfiguration()
	if err != nil {
		t.Errorf("#... - failed to parse configuration error: %s", err)
	}
	managerMock := &mocks.ResourceStatusInitiateManagerMock{}
	h := NewResourceStatusInitiateNotificationHandler(log, config, managerMock, requestName)
	return h, managerMock, config
}

func TestHandlerInit(t *testing.T) {
	h, _, _ := initRanConnectedNotificationHandlerTest(t, "RanConnected")
	assert.NotNil(t, h)
}

func TestHandleSuccess(t *testing.T) {
	h, managerMock, config := initRanConnectedNotificationHandlerTest(t, "RanConnected")

	payloadStr := "{\"nodeType\":1, \"messageDirection\":1}"
	payload := []byte(payloadStr)
	rmrReq := &models.RmrRequest{RanName:"RAN1", Payload:payload, Len:len(payload), StartTime:time.Now()}
	managerMock.On("Execute", rmrReq.RanName, mock.AnythingOfType("*e2pdus.ResourceStatusRequestData")).Return(nil)

	resourceStatusInitiateRequestParams := &e2pdus.ResourceStatusRequestData{}
	populateResourceStatusInitiateRequestParams(resourceStatusInitiateRequestParams, config)

	h.Handle(rmrReq)
	managerMock.AssertCalled(t, "Execute", rmrReq.RanName, resourceStatusInitiateRequestParams)
}

func TestHandleResourceStatusNotEnabled(t *testing.T) {
	h, managerMock, config := initRanConnectedNotificationHandlerTest(t, "RanConnected")
	config.ResourceStatusParams.EnableResourceStatus = false

	payloadStr := "{\"nodeType\":1, \"messageDirection\":1}"
	payload := []byte(payloadStr)
	rmrReq := &models.RmrRequest{RanName:"RAN1", Payload:payload, Len:len(payload), StartTime:time.Now()}
	managerMock.On("Execute", rmrReq.RanName, mock.AnythingOfType("*e2pdus.ResourceStatusRequestData")).Return(nil)

	h.Handle(rmrReq)
	managerMock.AssertNumberOfCalls(t, "Execute", 0)
}

func TestHandleUnknownJson(t *testing.T) {
	h, managerMock, _ := initRanConnectedNotificationHandlerTest(t, "RanConnected")

	payloadStr := "blablabla"
	payload := []byte(payloadStr)
	rmrReq := &models.RmrRequest{RanName:"RAN1", Payload:payload, Len:len(payload), StartTime:time.Now()}
	managerMock.On("Execute", rmrReq.RanName, mock.AnythingOfType("*e2pdus.ResourceStatusRequestData")).Return(nil)

	h.Handle(rmrReq)
	managerMock.AssertNumberOfCalls(t, "Execute", 0)
}

func TestHandleGnbNode(t *testing.T) {
	h, managerMock, _ := initRanConnectedNotificationHandlerTest(t, "RanConnected")

	payloadStr := "{\"nodeType\":2, \"messageDirection\":1}"
	payload := []byte(payloadStr)
	rmrReq := &models.RmrRequest{RanName:"RAN1", Payload:payload, Len:len(payload), StartTime:time.Now()}
	managerMock.On("Execute", rmrReq.RanName, mock.AnythingOfType("*e2pdus.ResourceStatusRequestData")).Return(nil)

	h.Handle(rmrReq)
	managerMock.AssertNumberOfCalls(t, "Execute", 0)
}
