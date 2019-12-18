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
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"github.com/pkg/errors"
	"rsm/configuration"
	"rsm/e2pdus"
	"rsm/enums"
	"rsm/logger"
	"rsm/mocks"
	"rsm/models"
	"rsm/services"
	"testing"
	"time"
)

func initResourceStatusResponseHandlerTest(t *testing.T) (*mocks.ResourceStatusResponseConverterMock, ResourceStatusResponseHandler, *mocks.RsmReaderMock, *mocks.RsmWriterMock) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}

	config, err := configuration.ParseConfiguration()
	if err != nil {
		t.Errorf("#... - failed to parse configuration error: %s", err)
	}

	converterMock := &mocks.ResourceStatusResponseConverterMock{}

	rnibReaderMock := &mocks.RnibReaderMock{}
	rsmReaderMock := &mocks.RsmReaderMock{}
	rsmWriterMock := &mocks.RsmWriterMock{}

	rnibDataService := services.NewRnibDataService(logger, config, rnibReaderMock, rsmReaderMock, rsmWriterMock)
	h := NewResourceStatusResponseHandler(logger, converterMock, rnibDataService)

	return converterMock, h, rsmReaderMock, rsmWriterMock

}

// Verify UnpackX2apPduAsString() and Convert() are called
func TestResourceStatusResponseHandler(t *testing.T) {
	converterMock, h, _, _ := initResourceStatusResponseHandlerTest(t)

	payload := []byte("aaa")
	req := models.RmrRequest{RanName: "test", StartTime: time.Now(), Payload: payload, Len: len(payload)}
	converterMock.On("UnpackX2apPduAsString", req.Payload, e2pdus.MaxAsn1CodecMessageBufferSize).Return(string(payload), nil)
	converterMock.On("Convert", req.Payload).Return((*models.ResourceStatusResponse)(nil), fmt.Errorf("error"))

	h.Handle(&req)
	converterMock.AssertNumberOfCalls(t, "UnpackX2apPduAsString", 1)
	converterMock.AssertNumberOfCalls(t, "Convert", 1)
}

func TestResourceStatusResponseHandlerConvertError(t *testing.T) {
	converterMock, h, _, _ := initResourceStatusResponseHandlerTest(t)

	payload := []byte("aaa")
	req := models.RmrRequest{RanName: "test", StartTime: time.Now(), Payload: payload, Len: len(payload)}
	err := fmt.Errorf("error")
	var payloadAsString string
	converterMock.On("UnpackX2apPduAsString", req.Payload, e2pdus.MaxAsn1CodecMessageBufferSize).Return(payloadAsString, err)
	converterMock.On("Convert", req.Payload).Return((*models.ResourceStatusResponse)(nil), fmt.Errorf("error"))
	h.Handle(&req)

	converterMock.AssertNumberOfCalls(t, "UnpackX2apPduAsString", 1)
	converterMock.AssertNumberOfCalls(t, "Convert", 0)
}

func TestResourceStatusResponseHandlerEnb2Mid0(t *testing.T) {
	converterMock, h, _, _ := initResourceStatusResponseHandlerTest(t)

	payload := []byte("aaa")
	req := models.RmrRequest{RanName: "test", StartTime: time.Now(), Payload: payload, Len: len(payload)}
	response := &models.ResourceStatusResponse{}
	converterMock.On("UnpackX2apPduAsString", req.Payload, e2pdus.MaxAsn1CodecMessageBufferSize).Return(string(payload), nil)
	converterMock.On("Convert", req.Payload).Return(response, nil)
	h.Handle(&req)

	converterMock.AssertNumberOfCalls(t, "UnpackX2apPduAsString", 1)
	converterMock.AssertNumberOfCalls(t, "Convert", 1)
}

func TestResourceStatusResponseHandlerWithMidGetRsmRanInfoFailure(t *testing.T) {
	converterMock, h, rsmReaderMock, rsmWriterMock := initResourceStatusResponseHandlerTest(t)
	payload := []byte("aaa")
	req := models.RmrRequest{RanName: RanName, StartTime: time.Now(), Payload: payload, Len: len(payload)}
	response := &models.ResourceStatusResponse{ENB1_Measurement_ID: 1, ENB2_Measurement_ID: 2}
	converterMock.On("UnpackX2apPduAsString", req.Payload, e2pdus.MaxAsn1CodecMessageBufferSize).Return(string(payload), nil)
	converterMock.On("Convert", req.Payload).Return(response, nil)
	rsmReaderMock.On("GetRsmRanInfo", RanName).Return(&models.RsmRanInfo{}, common.NewInternalError(errors.New("Error")))
	h.Handle(&req)
	converterMock.AssertNumberOfCalls(t, "UnpackX2apPduAsString", 1)
	converterMock.AssertNumberOfCalls(t, "Convert", 1)
	rsmReaderMock.AssertCalled(t, "GetRsmRanInfo", RanName)
	rsmWriterMock.AssertNotCalled(t, "SaveRsmRanInfo")
}

func TestResourceStatusResponseHandlerWithMidUpdateFailure(t *testing.T) {
	converterMock, h, rsmReaderMock, rsmWriterMock := initResourceStatusResponseHandlerTest(t)
	payload := []byte("aaa")
	req := models.RmrRequest{RanName: RanName, StartTime: time.Now(), Payload: payload, Len: len(payload)}
	response := &models.ResourceStatusResponse{ENB1_Measurement_ID: 1, ENB2_Measurement_ID: 2}
	converterMock.On("UnpackX2apPduAsString", req.Payload, e2pdus.MaxAsn1CodecMessageBufferSize).Return(string(payload), nil)
	converterMock.On("Convert", req.Payload).Return(response, nil)
	rsmRanInfoBefore := models.NewRsmRanInfo(RanName, enums.Enb1MeasurementId, 0, enums.Start, false)
	rsmReaderMock.On("GetRsmRanInfo", RanName).Return(rsmRanInfoBefore, nil)
	updatedRsmRanInfo := models.NewRsmRanInfo(RanName, enums.Enb1MeasurementId, response.ENB2_Measurement_ID, enums.Start, true)
	rsmWriterMock.On("SaveRsmRanInfo", rsmRanInfoBefore).Return(common.NewInternalError(errors.New("Error")))
	h.Handle(&req)
	converterMock.AssertNumberOfCalls(t, "UnpackX2apPduAsString", 1)
	converterMock.AssertNumberOfCalls(t, "Convert", 1)
	rsmReaderMock.AssertCalled(t, "GetRsmRanInfo", RanName)
	rsmWriterMock.AssertCalled(t, "SaveRsmRanInfo", updatedRsmRanInfo)
}

func TestResourceStatusResponseHandlerWithMidSuccessfulUpdate(t *testing.T) {
	converterMock, h, rsmReaderMock, rsmWriterMock := initResourceStatusResponseHandlerTest(t)
	payload := []byte("aaa")
	req := models.RmrRequest{RanName: RanName, StartTime: time.Now(), Payload: payload, Len: len(payload)}
	response := &models.ResourceStatusResponse{ENB1_Measurement_ID: 1, ENB2_Measurement_ID: 2}
	converterMock.On("UnpackX2apPduAsString", req.Payload, e2pdus.MaxAsn1CodecMessageBufferSize).Return(string(payload), nil)
	converterMock.On("Convert", req.Payload).Return(response, nil)
	rsmRanInfoBefore := models.NewRsmRanInfo(RanName, enums.Enb1MeasurementId, 0, enums.Start, false)
	rsmReaderMock.On("GetRsmRanInfo", RanName).Return(rsmRanInfoBefore, nil)
	updatedRsmRanInfo := models.NewRsmRanInfo(RanName, enums.Enb1MeasurementId, response.ENB2_Measurement_ID, enums.Start, true)
	rsmWriterMock.On("SaveRsmRanInfo", rsmRanInfoBefore).Return(nil)
	h.Handle(&req)
	converterMock.AssertNumberOfCalls(t, "UnpackX2apPduAsString", 1)
	converterMock.AssertNumberOfCalls(t, "Convert", 1)
	rsmReaderMock.AssertCalled(t, "GetRsmRanInfo", RanName)
	rsmWriterMock.AssertCalled(t, "SaveRsmRanInfo", updatedRsmRanInfo /*&updatedRsmRanInfo*/)
}
