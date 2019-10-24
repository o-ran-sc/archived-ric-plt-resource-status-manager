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

// Verify UnpackX2apPduAsString is called
func TestFailureHandler(t *testing.T) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	payload:= []byte("aaa")
	req:= models.RmrRequest{RanName: "test", StartTime:time.Now(), Payload:payload,Len:len(payload)}
	unpackerMock:=mocks.Asn1PduUnpackerMock{}
	unpackerMock.On("UnpackX2apPduAsString", req.Len, req.Payload, e2pdus.MaxAsn1CodecMessageBufferSize).Return(string(payload), nil)
	h:= NewResourceStatusFailureHandler(logger, &unpackerMock)

	h.Handle(&req)

	unpackerMock.AssertNumberOfCalls(t, "UnpackX2apPduAsString", 1)
}


func TestFailureHandlerError(t *testing.T) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	payload:= []byte("aaa")
	req:= models.RmrRequest{RanName: "test", StartTime:time.Now(), Payload:payload,Len:len(payload)}
	unpackerMock:=mocks.Asn1PduUnpackerMock{}

	err = fmt.Errorf("error")
	var payloadAsString string
	unpackerMock.On("UnpackX2apPduAsString", req.Len, req.Payload, e2pdus.MaxAsn1CodecMessageBufferSize).Return(payloadAsString, err)
	h:= NewResourceStatusFailureHandler(logger, &unpackerMock)

	h.Handle(&req)

	unpackerMock.AssertNumberOfCalls(t, "UnpackX2apPduAsString", 1)
}
