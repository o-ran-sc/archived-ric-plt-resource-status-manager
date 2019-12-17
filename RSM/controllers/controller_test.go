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


package controllers

import (
	"encoding/json"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"rsm/configuration"
	"rsm/enums"
	"rsm/logger"
	"rsm/mocks"
	"rsm/models"
	"rsm/providers/httpmsghandlerprovider"
	"rsm/rsmerrors"
	"rsm/services"
	"rsm/tests"
	"strings"
	"testing"
)

func setupControllerTest(t *testing.T) (*Controller, *mocks.RnibReaderMock, *mocks.RsmReaderMock, *mocks.RsmWriterMock, *mocks.ResourceStatusServiceMock) {
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
	handlerProvider := httpmsghandlerprovider.NewRequestHandlerProvider(log, rnibDataService, resourceStatusServiceMock)
	controller := NewController(log, handlerProvider)
	return controller, rnibReaderMock, rsmReaderMock, rsmWriterMock, resourceStatusServiceMock
}

func TestResourceStatusInvalidBody(t *testing.T) {
	controller, _, _, _ , _:= setupControllerTest(t)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	httpRequest, _ := http.NewRequest("PUT", "http://localhost:4800/v1/general/resourcestatus", strings.NewReader("{}{}"))
	httpRequest.Header = header

	writer := httptest.NewRecorder()
	controller.ResourceStatus(writer, httpRequest)

	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)
	assert.Equal(t, rsmerrors.NewInvalidJsonError().Code, errorResponse.Code)
}

func TestResourceStatusSuccess(t *testing.T) {
	controller, readerMock, rsmReaderMock, rsmWriterMock,  resourceStatusServiceMock := setupControllerTest(t)

	cellId1 := "02f829:0007ab00"
	cellId2 := "02f829:0007ab50"
	nodebInfo := &entities.NodebInfo{
		RanName:         tests.RanName,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Configuration: &entities.NodebInfo_Enb{
			Enb: &entities.Enb{
				ServedCells: []*entities.ServedCellInfo{{CellId: cellId1}, {CellId: cellId2}},
			},
		},
	}
	var nbIdentityList []*entities.NbIdentity
	config := tests.GetRsmGeneralConfiguration(true)

	rsmWriterMock.On("SaveRsmGeneralConfiguration", config).Return(nil)
	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(config, nil)
	readerMock.On("GetListEnbIds").Return(nbIdentityList, nil)
	readerMock.On("GetNodeb", tests.RanName).Return(nodebInfo)
	resourceStatusServiceMock.On("BuildAndSendInitiateRequest", mock.AnythingOfType("*entities.NodebInfo"), mock.AnythingOfType("*models.RsmGeneralConfiguration"), enums.Enb1MeasurementId).Return(nil)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	httpRequest, _ := http.NewRequest("PUT", "http://localhost:4800/v1/general/resourcestatus", strings.NewReader("{\"enableResourceStatus\":true}"))
	httpRequest.Header = header

	writer := httptest.NewRecorder()
	controller.ResourceStatus(writer, httpRequest)

	readerMock.AssertNumberOfCalls(t, "GetListEnbIds", 1)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 0)

	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)
}

func TestResourceStatusFail(t *testing.T) {

	controller, readerMock, rsmReaderMock, _,  resourceStatusServiceMock := setupControllerTest(t)

	rnibErr := &rsmerrors.RnibDbError{}

	rsmReaderMock.On("GetRsmGeneralConfiguration").Return(tests.GetRsmGeneralConfiguration(true), rnibErr)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	httpRequest, _ := http.NewRequest("PUT", "http://localhost:4800/v1/general/resourcestatus", strings.NewReader("{\"enableResourceStatus\":true}"))
	httpRequest.Header = header

	writer := httptest.NewRecorder()
	controller.ResourceStatus(writer, httpRequest)

	readerMock.AssertNumberOfCalls(t, "GetListEnbIds", 0)
	resourceStatusServiceMock.AssertNumberOfCalls(t, "BuildAndSendInitiateRequest", 0)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
}

func TestHeaderValidationFailed(t *testing.T) {
	controller, _, _ ,_ , _:= setupControllerTest(t)

	writer := httptest.NewRecorder()
	header := &http.Header{}
	controller.handleRequest(writer, header, "", nil, true)
	var errorResponse = parseJsonRequest(t, writer.Body)
	err := rsmerrors.NewHeaderValidationError()

	assert.Equal(t, http.StatusUnsupportedMediaType, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestHandleInternalError(t *testing.T) {
	controller, _, _,_ ,_ := setupControllerTest(t)

	writer := httptest.NewRecorder()
	err := rsmerrors.NewInternalError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestValidateHeadersSuccess(t *testing.T) {
	controller, _, _,_ ,_ := setupControllerTest(t)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	result := controller.validateRequestHeader(&header)

	assert.Nil(t, result)
}

func parseJsonRequest(t *testing.T, r io.Reader) models.ErrorResponse {

	var errorResponse models.ErrorResponse
	body, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Error cannot deserialize json request")
	}
	json.Unmarshal(body, &errorResponse)

	return errorResponse
}

func TestHandleErrorResponse(t *testing.T) {
	controller, _, _ ,_ , _:= setupControllerTest(t)

	writer := httptest.NewRecorder()
	controller.handleErrorResponse(rsmerrors.NewRnibDbError(), writer)
	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(rsmerrors.NewHeaderValidationError(), writer)
	assert.Equal(t, http.StatusUnsupportedMediaType, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(rsmerrors.NewWrongStateError("", ""), writer)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(rsmerrors.NewRequestValidationError(), writer)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(rsmerrors.NewRmrError(), writer)
	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(rsmerrors.NewResourceNotFoundError(), writer)
	assert.Equal(t, http.StatusNotFound, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(rsmerrors.NewInvalidJsonError(), writer)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(fmt.Errorf("ErrorError"), writer)
	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(rsmerrors.NewRsmError(2), writer)
	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
}