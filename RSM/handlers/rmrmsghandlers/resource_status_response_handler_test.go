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
	"rsm/e2pdus"
	"rsm/logger"
	"rsm/mocks"
	"rsm/models"
	"testing"
	"time"
)

// Verify UnpackX2apPduAsString() and Convert() are called
func TestResourceStatusResponseHandler(t *testing.T) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	payload := []byte("aaa")
	req := models.RmrRequest{RanName: "test", StartTime: time.Now(), Payload: payload, Len: len(payload)}
	converterMock := mocks.ResourceStatusResponseConverterMock{}
	converterMock.On("UnpackX2apPduAsString", req.Payload, e2pdus.MaxAsn1CodecMessageBufferSize).Return(string(payload), nil)
	converterMock.On("Convert", req.Payload).Return((*models.ResourceStatusResponse)(nil), fmt.Errorf("error"))
	h := NewResourceStatusResponseHandler(logger, &converterMock)

	h.Handle(&req)

	converterMock.AssertNumberOfCalls(t, "UnpackX2apPduAsString", 1)
	converterMock.AssertNumberOfCalls(t, "Convert", 1)
}

func TestResourceStatusResponseHandlerError(t *testing.T) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	payload := []byte("aaa")
	req := models.RmrRequest{RanName: "test", StartTime: time.Now(), Payload: payload, Len: len(payload)}
	converterMock := mocks.ResourceStatusResponseConverterMock{}

	err = fmt.Errorf("error")
	var payloadAsString string
	converterMock.On("UnpackX2apPduAsString", req.Payload, e2pdus.MaxAsn1CodecMessageBufferSize).Return(payloadAsString, err)
	converterMock.On("Convert", req.Payload).Return((*models.ResourceStatusResponse)(nil), fmt.Errorf("error"))
	h := NewResourceStatusResponseHandler(logger, &converterMock)

	h.Handle(&req)

	converterMock.AssertNumberOfCalls(t, "UnpackX2apPduAsString", 1)
	converterMock.AssertNumberOfCalls(t, "Convert", 0)
}

func TestResourceStatusResponseHandlerEnb2Mid0(t *testing.T) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	payload := []byte("aaa")
	req := models.RmrRequest{RanName: "test", StartTime: time.Now(), Payload: payload, Len: len(payload)}
	response := &models.ResourceStatusResponse{}
	converterMock := mocks.ResourceStatusResponseConverterMock{}
	converterMock.On("UnpackX2apPduAsString", req.Payload, e2pdus.MaxAsn1CodecMessageBufferSize).Return(string(payload), nil)
	converterMock.On("Convert", req.Payload).Return(response, nil)
	h := NewResourceStatusResponseHandler(logger, &converterMock)

	h.Handle(&req)

	converterMock.AssertNumberOfCalls(t, "UnpackX2apPduAsString", 1)
	converterMock.AssertNumberOfCalls(t, "Convert", 1)
}

func TestResourceStatusResponseHandlerWithMid(t *testing.T) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	payload := []byte("aaa")
	req := models.RmrRequest{RanName: "test", StartTime: time.Now(), Payload: payload, Len: len(payload)}
	response := &models.ResourceStatusResponse{ENB1_Measurement_ID: 1, ENB2_Measurement_ID: 2}
	converterMock := mocks.ResourceStatusResponseConverterMock{}
	converterMock.On("UnpackX2apPduAsString", req.Payload, e2pdus.MaxAsn1CodecMessageBufferSize).Return(string(payload), nil)
	converterMock.On("Convert", req.Payload).Return(response, nil)
	h := NewResourceStatusResponseHandler(logger, &converterMock)

	h.Handle(&req)

	converterMock.AssertNumberOfCalls(t, "UnpackX2apPduAsString", 1)
	converterMock.AssertNumberOfCalls(t, "Convert", 1)
}
