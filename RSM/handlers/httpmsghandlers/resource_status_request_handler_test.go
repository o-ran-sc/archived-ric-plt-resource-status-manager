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
package httpmsghandlers_test

import (
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"rsm/configuration"
	"rsm/enums"
	"rsm/handlers/httpmsghandlers"
	"rsm/logger"
	"rsm/mocks"
	"rsm/models"
	"rsm/rsmerrors"
	"rsm/services"
	"rsm/tests"
	"testing"
)

func initTest(t *testing.T) (*httpmsghandlers.ResourceStatusRequestHandler, *mocks.RnibReaderMock, *mocks.RsmReaderMock, *mocks.RsmWriterMock, *mocks.ResourceStatusServiceMock) {
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
	handler := httpmsghandlers.NewResourceStatusRequestHandler(log, rnibDataService, resourceStatusServiceMock)

	return handler, rnibReaderMock, rsmReaderMock, rsmWriterMock, resourceStatusServiceMock
}

func TestResourceStatusRequestHandlerGetConfigError(t *testing.T) {

	handler, _, rsmReaderMock, _,  _ := initTest(t)

	err := common.NewInternalError(errors.New("Error"))
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(tests.GetRsmGeneralConfiguration(true), err)

	resourceStatusRequest := models.ResourceStatusRequest{EnableResourceStatus:true}
	actualErr := handler.Handle(resourceStatusRequest)

	rsmReaderMock.AssertNumberOfCalls(t, "SaveRsmGeneralConfiguration", 0)

	assert.Equal(t, actualErr, rsmerrors.NewRnibDbError())
}

func TestResourceStatusRequestHandlerSaveConfigError(t *testing.T) {

	handler, readerMock, rsmReaderMock, rsmWriterMock,  _ := initTest(t)

	err := common.NewInternalError(errors.New("Error"))
	config := tests.GetRsmGeneralConfiguration(true)
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(config, nil)
	rsmWriterMock.On("SaveRsmGeneralConfiguration", config).Return(err)

	resourceStatusRequest := models.ResourceStatusRequest{EnableResourceStatus:true}
	actualErr := handler.Handle(resourceStatusRequest)

	readerMock.AssertNumberOfCalls(t, "GetListEnbIds", 0)

	assert.Equal(t, actualErr, rsmerrors.NewRnibDbError())
}

func TestResourceStatusRequestHandleGetListEnbIdsError(t *testing.T) {

	handler, readerMock, rsmReaderMock, rsmWriterMock,  _ := initTest(t)

	err := common.NewInternalError(errors.New("Error"))
	config := tests.GetRsmGeneralConfiguration(true)
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(config, nil)
	rsmWriterMock.On("SaveRsmGeneralConfiguration", config).Return(nil)

	var nbIdentityList []*entities.NbIdentity
	readerMock.On("GetListEnbIds").Return(nbIdentityList, err)

	resourceStatusRequest := models.ResourceStatusRequest{EnableResourceStatus:true}
	actualErr := handler.Handle(resourceStatusRequest)

	readerMock.AssertNumberOfCalls(t, "GetNodeb", 0)

	assert.Equal(t, actualErr, rsmerrors.NewRnibDbError())
}

func TestResourceStatusRequestHandlerTrueStartSuccess(t *testing.T) {

	handler, readerMock, rsmReaderMock, rsmWriterMock, resourceStatusServiceMock := initTest(t)

	config := tests.GetRsmGeneralConfiguration(true)
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(config, nil)
	rsmWriterMock.On("SaveRsmGeneralConfiguration", config).Return(nil)

	nbIdentityList := CreateIdentityList()
	readerMock.On("GetListEnbIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	nb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3, nil)

	rrInfo1 := &models.RsmRanInfo{RanName:"RanName_1", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:0, Action:enums.Start, ActionStatus:false}
	rrInfo2 := &models.RsmRanInfo{RanName:"RanName_2", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:0, Action:enums.Start, ActionStatus:true}
	rrInfo3 := &models.RsmRanInfo{RanName:"RanName_3", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:0, Action:enums.Stop, ActionStatus:false}
	rsmReaderMock.On("GetRsmRanInfo", "RanName_1").Return(rrInfo1, nil)
	rsmReaderMock.On("GetRsmRanInfo", "RanName_2").Return(rrInfo2, nil)
	rsmReaderMock.On("GetRsmRanInfo", "RanName_3").Return(rrInfo3, nil)

	rsmWriterMock.On("SaveRsmRanInfo", rrInfo3).Return(nil)

	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", nb1, config, enums.Enb1MeasurementId).Return(nil)
	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", nb3, config, enums.Enb1MeasurementId).Return(nil)


	resourceStatusRequest := models.ResourceStatusRequest{EnableResourceStatus:true}
	actualErr := handler.Handle(resourceStatusRequest)

	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	rsmWriterMock.AssertNumberOfCalls(t, "SaveRsmRanInfo", 1)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 2)

	assert.Equal(t, actualErr, nil)
}

func TestResourceStatusRequestHandlerTrueNumberOfFails2(t *testing.T) {

	handler, readerMock, rsmReaderMock, rsmWriterMock, resourceStatusServiceMock := initTest(t)

	config := tests.GetRsmGeneralConfiguration(true)
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(config, nil)
	rsmWriterMock.On("SaveRsmGeneralConfiguration", config).Return(nil)

	nbIdentityList := CreateIdentityList()
	readerMock.On("GetListEnbIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED,}
	nb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3, nil)

	err := common.NewInternalError(errors.New("Error"))
	rrInfo1 := &models.RsmRanInfo{RanName:"RanName_1", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:0, Action:enums.Start, ActionStatus:false}
	rrInfo3 := &models.RsmRanInfo{RanName:"RanName_3", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:0, Action:enums.Stop, ActionStatus:false}
	rsmReaderMock.On("GetRsmRanInfo", "RanName_1").Return(rrInfo1, err)
	rsmReaderMock.On("GetRsmRanInfo", "RanName_3").Return(rrInfo3, nil)

	rsmWriterMock.On("SaveRsmRanInfo", rrInfo3).Return(nil)

	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", nb3, mock.AnythingOfType("*models.RsmGeneralConfiguration"), enums.Enb1MeasurementId).Return(nil)


	resourceStatusRequest := models.ResourceStatusRequest{EnableResourceStatus:true}
	actualErr := handler.Handle(resourceStatusRequest)

	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	rsmWriterMock.AssertNumberOfCalls(t, "SaveRsmRanInfo", 1)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 1)

	rsmError := rsmerrors.NewRsmError(2)
	assert.Equal(t, actualErr, rsmError)
	assert.Equal(t, actualErr.Error(), rsmError.Error())
}

func TestResourceStatusRequestHandlerTrueNumberOfFails3(t *testing.T) {

	handler, readerMock, rsmReaderMock, rsmWriterMock, resourceStatusServiceMock := initTest(t)

	config := tests.GetRsmGeneralConfiguration(true)
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(config, nil)
	rsmWriterMock.On("SaveRsmGeneralConfiguration", config).Return(nil)

	nbIdentityList := CreateIdentityList()
	readerMock.On("GetListEnbIds").Return(nbIdentityList, nil)

	err := common.NewInternalError(errors.New("Error"))
	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	nb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, err)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3, nil)

	rrInfo1 := &models.RsmRanInfo{RanName:"RanName_1", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:0, Action:enums.Start, ActionStatus:false}
	rrInfo3 := &models.RsmRanInfo{RanName:"RanName_3", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:0, Action:enums.Stop, ActionStatus:false}
	rsmReaderMock.On("GetRsmRanInfo", "RanName_1").Return(rrInfo1, nil)
	rsmReaderMock.On("GetRsmRanInfo", "RanName_3").Return(rrInfo3, nil)

	rsmWriterMock.On("SaveRsmRanInfo", rrInfo3).Return(err)

	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", nb1, mock.AnythingOfType("*models.RsmGeneralConfiguration"), enums.Enb1MeasurementId).Return(errors.New("Error"))


	resourceStatusRequest := models.ResourceStatusRequest{EnableResourceStatus:true}
	actualErr := handler.Handle(resourceStatusRequest)

	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	rsmWriterMock.AssertNumberOfCalls(t, "SaveRsmRanInfo", 1)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 1)

	rsmError := rsmerrors.NewRsmError(3)
	assert.Equal(t, actualErr, rsmError)
	assert.Equal(t, actualErr.Error(), rsmError.Error())
}

func TestResourceStatusRequestHandlerFalseStopSuccess(t *testing.T) {

	handler, readerMock, rsmReaderMock, rsmWriterMock, resourceStatusServiceMock := initTest(t)

	config := tests.GetRsmGeneralConfiguration(true)
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(config, nil)
	rsmWriterMock.On("SaveRsmGeneralConfiguration", config).Return(nil)

	nbIdentityList := CreateIdentityList()
	readerMock.On("GetListEnbIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	nb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3, nil)

	rrInfo1 := &models.RsmRanInfo{RanName:"RanName_1", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:2, Action:enums.Stop, ActionStatus:false}
	rrInfo2 := &models.RsmRanInfo{RanName:"RanName_2", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:2, Action:enums.Stop, ActionStatus:true}
	rrInfo3 := &models.RsmRanInfo{RanName:"RanName_3", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:2, Action:enums.Start, ActionStatus:false}
	rsmReaderMock.On("GetRsmRanInfo", "RanName_1").Return(rrInfo1, nil)
	rsmReaderMock.On("GetRsmRanInfo", "RanName_2").Return(rrInfo2, nil)
	rsmReaderMock.On("GetRsmRanInfo", "RanName_3").Return(rrInfo3, nil)

	rsmWriterMock.On("SaveRsmRanInfo", rrInfo3).Return(nil)

	resourceStatusServiceMock.On("BuildAndSendStopRequest", nb1, config, rrInfo1.Enb1MeasurementId, rrInfo1.Enb2MeasurementId).Return(nil)
	resourceStatusServiceMock.On("BuildAndSendStopRequest", nb3, config, rrInfo3.Enb1MeasurementId, rrInfo3.Enb2MeasurementId).Return(nil)


	resourceStatusRequest := models.ResourceStatusRequest{EnableResourceStatus:false}
	actualErr := handler.Handle(resourceStatusRequest)

	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	rsmWriterMock.AssertNumberOfCalls(t, "SaveRsmRanInfo", 1)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendStopRequest", 2)

	assert.Equal(t, actualErr, nil)
}

func TestResourceStatusRequestHandlerFalseNumberOfFails2(t *testing.T) {

	handler, readerMock, rsmReaderMock, rsmWriterMock, resourceStatusServiceMock := initTest(t)

	config := tests.GetRsmGeneralConfiguration(true)
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(config, nil)
	rsmWriterMock.On("SaveRsmGeneralConfiguration", config).Return(nil)

	nbIdentityList := CreateIdentityList()
	readerMock.On("GetListEnbIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED,}
	nb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3, nil)

	err := common.NewInternalError(errors.New("Error"))
	rrInfo1 := &models.RsmRanInfo{RanName:"RanName_1", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:2, Action:enums.Stop, ActionStatus:false}
	rrInfo3 := &models.RsmRanInfo{RanName:"RanName_3", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:2, Action:enums.Start, ActionStatus:false}
	rsmReaderMock.On("GetRsmRanInfo", "RanName_1").Return(rrInfo1, err)
	rsmReaderMock.On("GetRsmRanInfo", "RanName_3").Return(rrInfo3, nil)

	rsmWriterMock.On("SaveRsmRanInfo", rrInfo3).Return(nil)

	resourceStatusServiceMock.On("BuildAndSendStopRequest", nb3, config, rrInfo3.Enb1MeasurementId, rrInfo3.Enb2MeasurementId).Return(nil)


	resourceStatusRequest := models.ResourceStatusRequest{EnableResourceStatus:false}
	actualErr := handler.Handle(resourceStatusRequest)

	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	rsmWriterMock.AssertNumberOfCalls(t, "SaveRsmRanInfo", 1)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendStopRequest", 1)

	rsmError := rsmerrors.NewRsmError(2)
	assert.Equal(t, actualErr, rsmError)
	assert.Equal(t, actualErr.Error(), rsmError.Error())
}

func TestResourceStatusRequestHandlerFalseNumberOfFails3(t *testing.T) {

	handler, readerMock, rsmReaderMock, rsmWriterMock, resourceStatusServiceMock := initTest(t)

	config := tests.GetRsmGeneralConfiguration(true)
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(config, nil)
	rsmWriterMock.On("SaveRsmGeneralConfiguration", config).Return(nil)

	nbIdentityList := CreateIdentityList()
	readerMock.On("GetListEnbIds").Return(nbIdentityList, nil)

	err := common.NewInternalError(errors.New("Error"))
	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	nb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, err)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3, nil)

	rrInfo1 := &models.RsmRanInfo{RanName:"RanName_1", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:2, Action:enums.Stop, ActionStatus:false}
	rrInfo3 := &models.RsmRanInfo{RanName:"RanName_3", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:2, Action:enums.Start, ActionStatus:false}
	rsmReaderMock.On("GetRsmRanInfo", "RanName_1").Return(rrInfo1, nil)
	rsmReaderMock.On("GetRsmRanInfo", "RanName_3").Return(rrInfo3, nil)

	rsmWriterMock.On("SaveRsmRanInfo", rrInfo3).Return(err)

	resourceStatusServiceMock.On("BuildAndSendStopRequest", nb1, config, rrInfo1.Enb1MeasurementId, rrInfo1.Enb2MeasurementId).Return(errors.New("Error"))

	resourceStatusRequest := models.ResourceStatusRequest{EnableResourceStatus:false}
	actualErr := handler.Handle(resourceStatusRequest)

	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	rsmWriterMock.AssertNumberOfCalls(t, "SaveRsmRanInfo", 1)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendStopRequest", 1)

	rsmError := rsmerrors.NewRsmError(3)
	assert.Equal(t, actualErr, rsmError)
	assert.Equal(t, actualErr.Error(), rsmError.Error())
}

func TestResourceStatusRequestHandlerFalseNoEnb2MeasurementId(t *testing.T) {

	handler, readerMock, rsmReaderMock, rsmWriterMock, resourceStatusServiceMock := initTest(t)

	config := tests.GetRsmGeneralConfiguration(true)
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(config, nil)
	rsmWriterMock.On("SaveRsmGeneralConfiguration", config).Return(nil)

	nbIdentityList := CreateIdentityList()
	readerMock.On("GetListEnbIds").Return(nbIdentityList, nil)

	nb1 := &entities.NodebInfo{RanName: "RanName_1", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	nb2 := &entities.NodebInfo{RanName: "RanName_2", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	nb3 := &entities.NodebInfo{RanName: "RanName_3", ConnectionStatus: entities.ConnectionStatus_CONNECTED,}
	readerMock.On("GetNodeb", "RanName_1").Return(nb1, nil)
	readerMock.On("GetNodeb", "RanName_2").Return(nb2, nil)
	readerMock.On("GetNodeb", "RanName_3").Return(nb3, nil)

	rrInfo1 := &models.RsmRanInfo{RanName:"RanName_1", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:2, Action:enums.Stop, ActionStatus:false}
	rrInfo2 := &models.RsmRanInfo{RanName:"RanName_2", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:2, Action:enums.Stop, ActionStatus:true}
	rrInfo3 := &models.RsmRanInfo{RanName:"RanName_3", Enb1MeasurementId:enums.Enb1MeasurementId, Enb2MeasurementId:0, Action:enums.Start, ActionStatus:false}
	rsmReaderMock.On("GetRsmRanInfo", "RanName_1").Return(rrInfo1, nil)
	rsmReaderMock.On("GetRsmRanInfo", "RanName_2").Return(rrInfo2, nil)
	rsmReaderMock.On("GetRsmRanInfo", "RanName_3").Return(rrInfo3, nil)

	resourceStatusServiceMock.On("BuildAndSendStopRequest", nb1, config, rrInfo1.Enb1MeasurementId, rrInfo1.Enb2MeasurementId).Return(nil)

	resourceStatusRequest := models.ResourceStatusRequest{EnableResourceStatus:false}
	actualErr := handler.Handle(resourceStatusRequest)

	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	rsmWriterMock.AssertNumberOfCalls(t, "SaveRsmRanInfo", 0)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendStopRequest", 1)

	rsmError := rsmerrors.NewRsmError(1)
	assert.Equal(t, actualErr, rsmError)
	assert.Equal(t, actualErr.Error(), rsmError.Error())
}

func CreateIdentityList() []*entities.NbIdentity {
	nbIdentity1 := entities.NbIdentity{InventoryName: "RanName_1"}
	nbIdentity2 := entities.NbIdentity{InventoryName: "RanName_2"}
	nbIdentity3 := entities.NbIdentity{InventoryName: "RanName_3"}

	var nbIdentityList []*entities.NbIdentity
	nbIdentityList = append(nbIdentityList, &nbIdentity1)
	nbIdentityList = append(nbIdentityList, &nbIdentity2)
	nbIdentityList = append(nbIdentityList, &nbIdentity3)

	return nbIdentityList
}