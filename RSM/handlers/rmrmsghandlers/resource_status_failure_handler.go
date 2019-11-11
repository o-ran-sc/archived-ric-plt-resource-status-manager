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

type ResourceStatusFailureHandler struct {
	logger   *logger.Logger
	converter converters.IResourceStatusResponseConverter
}

func NewResourceStatusFailureHandler(logger *logger.Logger, converter converters.IResourceStatusResponseConverter) ResourceStatusFailureHandler {
	return ResourceStatusFailureHandler{
		logger:   logger,
		converter: converter,
	}
}

func (h ResourceStatusFailureHandler) Handle(request *models.RmrRequest) {
	h.logger.Infof("#ResourceStatusFailureHandler.Handle - RAN name: %s - Received resource status failure notification", request.RanName)
	if h.logger.DebugEnabled() {
		pduAsString, err := h.converter.UnpackX2apPduAsString(request.Len, request.Payload, e2pdus.MaxAsn1CodecMessageBufferSize)
		if err != nil {
			h.logger.Errorf("#ResourceStatusFailureHandler.Handle - RAN name: %s - unpack failed. Error: %v", request.RanName, err)
			return
		}
		h.logger.Infof("#ResourceStatusFailureHandler.Handle - RAN name: %s - message: %s", request.RanName, pduAsString)
	}
	response, err := h.converter.Convert(request.Len, request.Payload, e2pdus.MaxAsn1CodecMessageBufferSize)
	if err != nil {
		h.logger.Errorf("#ResourceStatusFailureHandler.Handle - RAN name: %s - unpack failed. Error: %v", request.RanName, err)
		return
	}

	/*The RSM creates one measurement per cell*/

	if response.MeasurementInitiationResults == nil {
		h.logger.Infof("#ResourceStatusFailureHandler.Handle - RAN name: %s - ENB1_Measurement_ID: %d, ENB2_Measurement_ID: %d",
			request.RanName,
			response.ENB1_Measurement_ID,
			response.ENB2_Measurement_ID)
	} else {
		h.logger.Infof("#ResourceStatusFailureHandler.Handle - RAN name: %s - ENB1_Measurement_ID: %d, ENB2_Measurement_ID: %d, CellId: %s, FailedReportCharacteristics: %x",
			request.RanName,
			response.ENB1_Measurement_ID,
			response.ENB2_Measurement_ID,
			response.MeasurementInitiationResults[0].CellId,
			response.MeasurementInitiationResults[0].MeasurementFailureCauses[0].MeasurementFailedReportCharacteristics)
	}
}
