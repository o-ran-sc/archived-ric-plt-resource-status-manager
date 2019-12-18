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

// Verify UnpackX2apPduAsString() and Convert() are called
/*
func TestResourceStatusFailureHandlerConvertFailure(t *testing.T) {
	logger, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	payload := []byte("aaa")
	req := models.RmrRequest{RanName: "test", StartTime: time.Now(), Payload: payload, Len: len(payload)}
	converterMock := mocks.ResourceStatusFailureConverterMock{}
	//converterMock.On("UnpackX2apPduAsString", req.Payload, e2pdus.MaxAsn1CodecMessageBufferSize).Return(string(payload), nil)
	converterMock.On("Convert", req.Payload).Return((*models.ResourceStatusResponse)(nil), fmt.Errorf("error"))
	h := NewResourceStatusFailureHandler(logger, &converterMock)

	h.Handle(&req)

	//converterMock.AssertNumberOfCalls(t, "UnpackX2apPduAsString", 1)
	converterMock.AssertNumberOfCalls(t, "Convert", 1)
}


func TestResourceStatusFailureHandlerUnpackFailure(t *testing.T) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	payload := []byte("aaa")
	req := models.RmrRequest{RanName: "test", StartTime: time.Now(), Payload: payload, Len: len(payload)}
	converterMock := mocks.ResourceStatusFailureConverterMock{}

	err = fmt.Errorf("error")
	var payloadAsString string
	converterMock.On("UnpackX2apPduAsString", req.Payload, e2pdus.MaxAsn1CodecMessageBufferSize).Return(payloadAsString, err)
	converterMock.On("Convert", req.Payload).Return((*models.ResourceStatusResponse)(nil), fmt.Errorf("error"))
	h := NewResourceStatusFailureHandler(logger, &converterMock)

	h.Handle(&req)

	converterMock.AssertNumberOfCalls(t, "UnpackX2apPduAsString", 1)
	converterMock.AssertNumberOfCalls(t, "Convert", 0)
}
*/

/*
func TestResourceStatusFailureHandler(t *testing.T) {
	logger, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	unpacker := converters.NewX2apPduUnpacker(logger, e2pdus.MaxAsn1CodecMessageBufferSize)
	converter := converters.NewResourceStatusFailureConverter(unpacker)
	var payload []byte
	fmt.Sscanf("400900320000040027000300000e0028000300000c00054001620044401800004540130002f8290007ab500000434006000000000740", "%x", &payload)
	req := models.RmrRequest{RanName: "test", StartTime: time.Now(), Payload: payload, Len: len(payload)}
	h := NewResourceStatusFailureHandler(logger, converter)

	h.Handle(&req)
}

func TestResourceStatusFailureHandlerMinimalPdu(t *testing.T) {
	logger, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	unpacker := converters.NewX2apPduUnpacker(logger, e2pdus.MaxAsn1CodecMessageBufferSize)
	converter := converters.NewResourceStatusFailureConverter(unpacker)
	var payload []byte
	fmt.Sscanf("400900170000030027000300000000280003000049000540020a80", "%x", &payload)
	req := models.RmrRequest{RanName: "test", StartTime: time.Now(), Payload: payload, Len: len(payload)}
	h := NewResourceStatusFailureHandler(logger, converter)

	h.Handle(&req)
}
*/