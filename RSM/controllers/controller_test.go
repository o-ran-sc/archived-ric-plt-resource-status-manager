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

package controllers

/*
import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"rsm/configuration"
	"rsm/mocks"
	"rsm/models"
	"rsm/providers/httpmsghandlerprovider"
	"rsm/rsmerrors"
	"rsm/tests/testhelper"
	"testing"
)

func setupControllerTest(t *testing.T) (*Controller, *mocks.RnibReaderMock, *mocks.RmrMessengerMock) {
	config := configuration.ParseConfiguration()

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	readerMock := &mocks.RnibReaderMock{}
	rnibDataService, rmrSender, log := testhelper.InitTestCase(t)
	handlerProvider := httpmsghandlerprovider.NewRequestHandlerProvider(log, rmrSender, config, rnibDataService)
	controller := NewController(log, handlerProvider)
	return controller, readerMock, rmrMessengerMock
}

func TestHeaderValidationFailed(t *testing.T) {
	controller, _, _ := setupControllerTest(t)

	writer := httptest.NewRecorder()
	header := &http.Header{}
	controller.handleRequest(writer, header, "", nil, true)
	var errorResponse = parseJsonRequest(t, writer.Body)
	err := rsmerrors.NewHeaderValidationError()

	assert.Equal(t, http.StatusUnsupportedMediaType, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}
*/
/*func TestHandleInternalError(t *testing.T) {
	controller, _, _ := setupControllerTest(t)

	writer := httptest.NewRecorder()
	err := rsmerrors.NewInternalError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}*/

/*func TestValidateHeadersSuccess(t *testing.T) {
	controller, _, _ := setupControllerTest(t)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	result := controller.validateRequestHeader(&header)

	assert.Nil(t, result)
}
*/
/*func parseJsonRequest(t *testing.T, r io.Reader) models.ErrorResponse {

	var errorResponse models.ErrorResponse
	body, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Error cannot deserialize json request")
	}
	json.Unmarshal(body, &errorResponse)

	return errorResponse
}*/
/*
func TestHandleErrorResponse(t *testing.T) {
	controller, _, _ := setupControllerTest(t)

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
}
*/