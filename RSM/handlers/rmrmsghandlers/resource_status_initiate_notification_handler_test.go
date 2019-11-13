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
package rmrmsghandlers

import (
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"rsm/configuration"
	"rsm/enums"
	"rsm/logger"
	"rsm/mocks"
	"rsm/models"
	"rsm/services"
	"testing"
	"time"
)

const RanName = "test"

func getRsmGeneralConfiguration(enableResourceStatus bool) *models.RsmGeneralConfiguration {
	return &models.RsmGeneralConfiguration{
		EnableResourceStatus:         enableResourceStatus,
		PartialSuccessAllowed:        true,
		PrbPeriodic:                  true,
		TnlLoadIndPeriodic:           true,
		HwLoadIndPeriodic:            true,
		AbsStatusPeriodic:            true,
		RsrpMeasurementPeriodic:      true,
		CsiPeriodic:                  true,
		PeriodicityMs:                enums.ReportingPeriodicity_one_thousand_ms,
		PeriodicityRsrpMeasurementMs: enums.ReportingPeriodicityRSRPMR_four_hundred_80_ms,
		PeriodicityCsiMs:             enums.ReportingPeriodicityCSIR_ms20,
	}
}

func initRanConnectedNotificationHandlerTest(t *testing.T, requestName string) (ResourceStatusInitiateNotificationHandler, *mocks.RnibReaderMock, *mocks.ResourceStatusServiceMock, *mocks.RsmWriterMock, *mocks.RsmReaderMock) {
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

	h := NewResourceStatusInitiateNotificationHandler(log, rnibDataService, resourceStatusServiceMock, requestName)
	return h, rnibReaderMock, resourceStatusServiceMock, rsmWriterMock, rsmReaderMock
}

func TestHandlerInit(t *testing.T) {
	h, _, _, _, _ := initRanConnectedNotificationHandlerTest(t, "RanConnected")
	assert.NotNil(t, h)
}

func TestJsonUnmarshalError(t *testing.T) {
	h, rnibReaderMock, resourceStatusServiceMock, _, _ := initRanConnectedNotificationHandlerTest(t, "RanConnected")

	payloadStr := "blabla"
	payload := []byte(payloadStr)
	rmrReq := &models.RmrRequest{RanName: RanName, Payload: payload, Len: len(payload), StartTime: time.Now()}

	rnibReaderMock.On("GetNodeb", RanName).Return(mock.AnythingOfType("*entities.NodebInfo"))
	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", mock.AnythingOfType("*entities.NodebInfo"), mock.AnythingOfType("*models.RsmGeneralConfiguration"), enums.Enb1MeasurementId).Return(nil)
	h.Handle(rmrReq)
	rnibReaderMock.AssertNumberOfCalls(t, "GetNodeb", 0)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 0)
}

func TestUnknownJsonValue(t *testing.T) {
	h, rnibReaderMock, resourceStatusServiceMock, _, _ := initRanConnectedNotificationHandlerTest(t, "RanConnected")

	payloadStr := "{\"whatever\":3}"
	payload := []byte(payloadStr)
	rmrReq := &models.RmrRequest{RanName: RanName, Payload: payload, Len: len(payload), StartTime: time.Now()}

	rnibReaderMock.On("GetNodeb", RanName).Return(mock.AnythingOfType("*entities.NodebInfo"))
	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", mock.AnythingOfType("*entities.NodebInfo"), mock.AnythingOfType("*models.RsmGeneralConfiguration"), enums.Enb1MeasurementId).Return(nil)
	h.Handle(rmrReq)
	rnibReaderMock.AssertNumberOfCalls(t, "GetNodeb", 0)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 0)
}

func TestHandleGnbNode(t *testing.T) {
	h, rnibReaderMock, resourceStatusServiceMock, _, _ := initRanConnectedNotificationHandlerTest(t, "RanConnected")

	payloadStr := "{\"nodeType\":2, \"messageDirection\":1}"
	payload := []byte(payloadStr)
	rmrReq := &models.RmrRequest{RanName: "RAN1", Payload: payload, Len: len(payload), StartTime: time.Now()}
	rnibReaderMock.On("GetNodeb", RanName).Return(mock.AnythingOfType("*entities.NodebInfo"))
	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", mock.AnythingOfType("*entities.NodebInfo"), mock.AnythingOfType("*models.RsmGeneralConfiguration"), enums.Enb1MeasurementId).Return(nil)
	h.Handle(rmrReq)
	rnibReaderMock.AssertNumberOfCalls(t, "GetNodeb", 0)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 0)
}

func TestGetNodebFailure(t *testing.T) {
	h, rnibReaderMock, resourceStatusServiceMock, _, _ := initRanConnectedNotificationHandlerTest(t, "RanConnected")
	var nodebInfo *entities.NodebInfo

	payloadStr := "{\"nodeType\":1, \"messageDirection\":1}"
	payload := []byte(payloadStr)
	rmrReq := &models.RmrRequest{RanName: RanName, Payload: payload, Len: len(payload), StartTime: time.Now()}

	rnibReaderMock.On("GetNodeb", RanName).Return(nodebInfo, common.NewInternalError(errors.New("Error")))
	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", mock.AnythingOfType("*entities.NodebInfo"), mock.AnythingOfType("*models.RsmGeneralConfiguration"), enums.Enb1MeasurementId).Return(nil)
	h.Handle(rmrReq)
	rnibReaderMock.AssertCalled(t, "GetNodeb", RanName)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 0)
}

func TestInvalidConnectionStatus(t *testing.T) {
	h, rnibReaderMock, resourceStatusServiceMock, _/*rsmWriterMock*/, _ := initRanConnectedNotificationHandlerTest(t, "RanConnected")
	var err error
	rnibReaderMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}, err)
	//rsmRanInfo := models.RsmRanInfo{RanName, 0, 0, enums.Stop, true}
	//rsmWriterMock.On("SaveRsmRanInfo", &rsmRanInfo).Return(err)

	payloadStr := "{\"nodeType\":1, \"messageDirection\":1}"
	payload := []byte(payloadStr)
	rmrReq := &models.RmrRequest{RanName: RanName, Payload: payload, Len: len(payload), StartTime: time.Now()}
	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", mock.AnythingOfType("*entities.NodebInfo"), mock.AnythingOfType("*models.RsmGeneralConfiguration"), enums.Enb1MeasurementId).Return(nil)
	h.Handle(rmrReq)
	rnibReaderMock.AssertCalled(t, "GetNodeb", RanName)
	//rsmWriterMock.AssertNumberOfCalls(t, "SaveRsmRanInfo", 1)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 0)
}

func TestGetRsmGeneralConfigurationFailure(t *testing.T) {
	h, rnibReaderMock, resourceStatusServiceMock, _, rsmReaderMock := initRanConnectedNotificationHandlerTest(t, "RanConnected")
	var err error
	rnibReaderMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}, err)
	var rgc models.RsmGeneralConfiguration
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(&rgc, common.NewInternalError(errors.New("Error")))
	payloadStr := "{\"nodeType\":1, \"messageDirection\":1}"
	payload := []byte(payloadStr)
	rmrReq := &models.RmrRequest{RanName: RanName, Payload: payload, Len: len(payload), StartTime: time.Now()}
	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", mock.AnythingOfType("*entities.NodebInfo"), mock.AnythingOfType("*models.RsmGeneralConfiguration"), enums.Enb1MeasurementId).Return(nil)
	h.Handle(rmrReq)
	rnibReaderMock.AssertCalled(t, "GetNodeb", RanName)
	rsmReaderMock.AssertCalled(t, "GetRsmGeneralConfiguration")
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 0)
}

func TestEnableResourceStatusFalse(t *testing.T) {
	h, rnibReaderMock, resourceStatusServiceMock, rsmWriterMock, rsmReaderMock := initRanConnectedNotificationHandlerTest(t, "RanConnected")
	var err error
	rnibReaderMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}, err)
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(getRsmGeneralConfiguration(false), err)
	rsmRanInfo := models.RsmRanInfo{RanName, 0, 0, enums.Stop, true}
	rsmWriterMock.On("SaveRsmRanInfo", &rsmRanInfo).Return(err)

	payloadStr := "{\"nodeType\":1, \"messageDirection\":1}"
	payload := []byte(payloadStr)
	rmrReq := &models.RmrRequest{RanName: RanName, Payload: payload, Len: len(payload), StartTime: time.Now()}
	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", mock.AnythingOfType("*entities.NodebInfo"), mock.AnythingOfType("*models.RsmGeneralConfiguration"), enums.Enb1MeasurementId).Return(nil)
	h.Handle(rmrReq)
	rnibReaderMock.AssertCalled(t, "GetNodeb", RanName)
	rsmReaderMock.AssertCalled(t, "GetRsmGeneralConfiguration")
	rsmWriterMock.AssertNumberOfCalls(t, "SaveRsmRanInfo", 1)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 0)
}

func TestEnableResourceStatusFalseRsmRanInfoFailure(t *testing.T) {
	h, rnibReaderMock, resourceStatusServiceMock, rsmWriterMock, rsmReaderMock := initRanConnectedNotificationHandlerTest(t, "RanConnected")
	var err error
	rnibReaderMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}, err)
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(getRsmGeneralConfiguration(false), err)
	rsmRanInfo := models.RsmRanInfo{RanName, 0, 0, enums.Stop, true}
	rsmWriterMock.On("SaveRsmRanInfo", &rsmRanInfo).Return(common.NewInternalError(errors.New("Error")))

	payloadStr := "{\"nodeType\":1, \"messageDirection\":1}"
	payload := []byte(payloadStr)
	rmrReq := &models.RmrRequest{RanName: RanName, Payload: payload, Len: len(payload), StartTime: time.Now()}
	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", mock.AnythingOfType("*entities.NodebInfo"), mock.AnythingOfType("*models.RsmGeneralConfiguration"), enums.Enb1MeasurementId).Return(nil)
	h.Handle(rmrReq)
	rnibReaderMock.AssertCalled(t, "GetNodeb", RanName)
	rsmReaderMock.AssertCalled(t, "GetRsmGeneralConfiguration")
	rsmWriterMock.AssertNumberOfCalls(t, "SaveRsmRanInfo", 1)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 0)
}


func TestEnableResourceStatusTrueSaveRsmRanInfoFailure(t *testing.T) {
	h, rnibReaderMock, resourceStatusServiceMock, rsmWriterMock, rsmReaderMock := initRanConnectedNotificationHandlerTest(t, "RanConnected")
	var err error
	rnibReaderMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}, err)
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(getRsmGeneralConfiguration(true), err)
	rsmRanInfo := models.RsmRanInfo{RanName, enums.Enb1MeasurementId, 0, enums.Start, false}
	rsmWriterMock.On("SaveRsmRanInfo", &rsmRanInfo).Return(common.NewInternalError(errors.New("Error")))

	payloadStr := "{\"nodeType\":1, \"messageDirection\":1}"
	payload := []byte(payloadStr)
	rmrReq := &models.RmrRequest{RanName: RanName, Payload: payload, Len: len(payload), StartTime: time.Now()}
	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", mock.AnythingOfType("*entities.NodebInfo"), mock.AnythingOfType("*models.RsmGeneralConfiguration"), enums.Enb1MeasurementId).Return(nil)
	h.Handle(rmrReq)
	rnibReaderMock.AssertCalled(t, "GetNodeb", RanName)
	rsmReaderMock.AssertCalled(t, "GetRsmGeneralConfiguration")
	rsmWriterMock.AssertNumberOfCalls(t, "SaveRsmRanInfo", 1)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 0)
}

func TestBuildAndSendSuccess(t *testing.T) {
	h, rnibReaderMock, resourceStatusServiceMock, rsmWriterMock, rsmReaderMock := initRanConnectedNotificationHandlerTest(t, "RanConnected")
	var err error
	nodebInfo := &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	rgc := getRsmGeneralConfiguration(true)
	rnibReaderMock.On("GetNodeb", RanName).Return(nodebInfo, err)
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(rgc, err)
	rsmRanInfo := models.RsmRanInfo{RanName, enums.Enb1MeasurementId, 0, enums.Start, false}
	rsmWriterMock.On("SaveRsmRanInfo", &rsmRanInfo).Return(err)

	payloadStr := "{\"nodeType\":1, \"messageDirection\":1}"
	payload := []byte(payloadStr)
	rmrReq := &models.RmrRequest{RanName: RanName, Payload: payload, Len: len(payload), StartTime: time.Now()}
	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", nodebInfo, rgc, enums.Enb1MeasurementId).Return(nil)
	h.Handle(rmrReq)
	rnibReaderMock.AssertCalled(t, "GetNodeb", RanName)
	rsmReaderMock.AssertCalled(t, "GetRsmGeneralConfiguration")
	rsmWriterMock.AssertNumberOfCalls(t, "SaveRsmRanInfo", 1)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 1)
}

func TestBuildAndSendError(t *testing.T) {
	h, rnibReaderMock, resourceStatusServiceMock, rsmWriterMock, rsmReaderMock := initRanConnectedNotificationHandlerTest(t, "RanConnected")
	var err error
	nodebInfo := &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	rgc := getRsmGeneralConfiguration(true)
	rnibReaderMock.On("GetNodeb", RanName).Return(nodebInfo, err)
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(rgc, err)
	rsmRanInfoStart := models.RsmRanInfo{RanName, enums.Enb1MeasurementId, 0, enums.Start, false}
	rsmWriterMock.On("SaveRsmRanInfo", &rsmRanInfoStart).Return(err)
	rsmRanInfoStop := models.RsmRanInfo{RanName, 0, 0, enums.Stop, true}
	rsmWriterMock.On("SaveRsmRanInfo", &rsmRanInfoStop).Return(err)
	payloadStr := "{\"nodeType\":1, \"messageDirection\":1}"
	payload := []byte(payloadStr)
	rmrReq := &models.RmrRequest{RanName: RanName, Payload: payload, Len: len(payload), StartTime: time.Now()}
	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", nodebInfo, rgc, enums.Enb1MeasurementId).Return(common.NewInternalError(errors.New("Error")))
	h.Handle(rmrReq)
	rnibReaderMock.AssertCalled(t, "GetNodeb", RanName)
	rsmReaderMock.AssertCalled(t, "GetRsmGeneralConfiguration")
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 1)
	rsmWriterMock.AssertNumberOfCalls(t, "SaveRsmRanInfo", 2)
}
