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
	"rsm/converters"
	"rsm/e2pdus"
	"rsm/logger"
	"rsm/models"
)

type ResourceStatusResponseHandler struct {
	logger *logger.Logger
	converter converters.IResourceStatusResponseConverter
}

func NewResourceStatusResponseHandler(logger *logger.Logger, converter converters.IResourceStatusResponseConverter) ResourceStatusResponseHandler {
	return ResourceStatusResponseHandler{
		logger: logger,
		converter: converter,
	}
}

func (h ResourceStatusResponseHandler) Handle(request *models.RmrRequest) {
	h.logger.Infof("#ResourceStatusResponseHandler.Handle - RAN name: %s - Received resource status response notification", request.RanName)
	if h.logger.DebugEnabled() {
		pduAsString, err := h.converter.UnpackX2apPduAsString(request.Len, request.Payload, e2pdus.MaxAsn1CodecMessageBufferSize)
		if err != nil {
			h.logger.Errorf("#ResourceStatusResponseHandler.Handle - RAN name: %s - unpack failed. Error: %v", request.RanName, err)
			return
		}
		h.logger.Debugf("#ResourceStatusResponseHandler.Handle - RAN name: %s - pdu: %s", request.RanName, pduAsString)
	}
	response, err := h.converter.Convert(request.Len, request.Payload, e2pdus.MaxAsn1CodecMessageBufferSize)
	if err != nil {
		h.logger.Errorf("#ResourceStatusResponseHandler.Handle - RAN name: %s - unpack failed. Error: %v", request.RanName, err)
		return
	}
	if response.ENB2_Measurement_ID == 0 {
		h.logger.Errorf("#ResourceStatusResponseHandler.Handle - RAN name: %s - ignoring response without ENB2_Measurement_ID for ENB1_Measurement_ID = %d", request.RanName, response.ENB1_Measurement_ID)
		return
	}


	h.logger.Infof("#ResourceStatusResponseHandler.Handle - RAN name: %s - (success) ENB1_Measurement_ID: %d, ENB2_Measurement_ID: %d",
		request.RanName,
		response.ENB1_Measurement_ID,
		response.ENB2_Measurement_ID)
}
