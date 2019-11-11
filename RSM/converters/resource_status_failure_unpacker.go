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

package converters

// #cgo CFLAGS: -I../asn1codec/inc/ -I../asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../asn1codec/lib/ -L../asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <asn1codec_utils.h>
// #include <resource_status_response_wrapper.h>
import "C"
import (
	"fmt"
	"rsm/models"
	"unsafe"
)

type ResourceStatusFailureConverter struct {
	X2apPduUnpacker
}

func NewResourceStatusFailureConverter(unpacker X2apPduUnpacker) ResourceStatusFailureConverter {
	return ResourceStatusFailureConverter{unpacker}
}

func convertCompleteFailureCauseInformation(completeFailureCauseInformation_List *C.CompleteFailureCauseInformation_List_t) ([]*models.MeasurementInitiationResult, error) {
	var measurementInitiationResults []*models.MeasurementInitiationResult

	count := int(completeFailureCauseInformation_List.list.count)
	if count < 1 || count > maxCellineNB {
		return nil, fmt.Errorf("invalid number of complete failure cause information elements, %d", count)
	}
	completeFailureCauseInformation_ItemIEs_slice := (*[1 << 30]*C.CompleteFailureCauseInformation_ItemIEs_t)(unsafe.Pointer(completeFailureCauseInformation_List.list.array))[:count:count]
	for _, itemIE := range completeFailureCauseInformation_ItemIEs_slice {

		switch itemIE.value.present {
		case C.CompleteFailureCauseInformation_ItemIEs__value_PR_CompleteFailureCauseInformation_Item:
			item := (*C.CompleteFailureCauseInformation_Item_t)(unsafe.Pointer(&itemIE.value.choice[0]))
			measurementInitiationResult := models.MeasurementInitiationResult{CellId: buildCellId(item.cell_ID)}
			convertMeasurementFailureCauses(&item.measurementFailureCause_List, &measurementInitiationResult)
			measurementInitiationResults = append(measurementInitiationResults, &measurementInitiationResult)
		}
	}

	return measurementInitiationResults, nil
}

func convertResourceFailureIEs(resourceStatusFailure *C.ResourceStatusFailure_t, response *models.ResourceStatusResponse) error {
	count := int(resourceStatusFailure.protocolIEs.list.count)
	resourceStatusFailure_IEs_slice := (*[1 << 30]*C.ResourceStatusFailure_IEs_t)(unsafe.Pointer(resourceStatusFailure.protocolIEs.list.array))[:count:count]
	for _, resourceStatusFailure_IEs := range resourceStatusFailure_IEs_slice {
		switch resourceStatusFailure_IEs.value.present {
		case C.ResourceStatusFailure_IEs__value_PR_Measurement_ID:
			measurement_ID := (*C.Measurement_ID_t)(unsafe.Pointer(&resourceStatusFailure_IEs.value.choice[0]))
			if resourceStatusFailure_IEs.id == C.ProtocolIE_ID_id_ENB1_Measurement_ID {
				response.ENB1_Measurement_ID = int64(*measurement_ID)
			}
			if resourceStatusFailure_IEs.id == C.ProtocolIE_ID_id_ENB2_Measurement_ID {
				response.ENB2_Measurement_ID = int64(*measurement_ID)
			}
		case C.ResourceStatusFailure_IEs__value_PR_CriticalityDiagnostics:
			/*ignored*/
		case C.ResourceStatusFailure_IEs__value_PR_Cause:
			/*ignored*/
		case C.ResourceStatusFailure_IEs__value_PR_CompleteFailureCauseInformation_List:
			measurementInitiationResults, err := convertCompleteFailureCauseInformation((*C.CompleteFailureCauseInformation_List_t)(unsafe.Pointer(&resourceStatusFailure_IEs.value.choice[0])))
			if err != nil {
				return err
			}
			response.MeasurementInitiationResults = measurementInitiationResults
		}
	}
	return nil
}

func convertResourceStatusFailure(pdu *C.E2AP_PDU_t) (*models.ResourceStatusResponse, error) {
	response := models.ResourceStatusResponse{}

	if pdu.present != C.E2AP_PDU_PR_unsuccessfulOutcome {
		return &response, fmt.Errorf("unexpected PDU, %d", int(pdu.present))
	}

	//dereference a union of pointers (C union is represented as a byte array with the size of the largest member)
	unsuccessfulOutcome := *(**C.UnsuccessfulOutcome_t)(unsafe.Pointer(&pdu.choice[0]))
	if unsuccessfulOutcome == nil || unsuccessfulOutcome.value.present != C.UnsuccessfulOutcome__value_PR_ResourceStatusFailure {
		return &response, fmt.Errorf("unexpected PDU - not a resource status failure")
	}

	resourceStatusFailure := (*C.ResourceStatusFailure_t)(unsafe.Pointer(&unsuccessfulOutcome.value.choice[0]))
	if resourceStatusFailure.protocolIEs.list.count == 0 {
		return &response, fmt.Errorf("unexpected PDU - no protocolIEs found")
	}

	if err := convertResourceFailureIEs(resourceStatusFailure, &response); err != nil {
		return &response, err
	}

	return &response, nil
}

// Convert pdu to public ResourceStatusResponse
func (r ResourceStatusFailureConverter) Convert(packedBufferSize int, packedBuf []byte, maxMessageBufferSize int) (*models.ResourceStatusResponse, error) {
	pdu, err := r.UnpackX2apPdu(packedBufferSize, packedBuf, maxMessageBufferSize)
	if err != nil {
		return nil, err
	}

	defer C.delete_pdu(pdu)
	return convertResourceStatusFailure(pdu)
}
