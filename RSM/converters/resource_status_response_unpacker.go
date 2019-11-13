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

const (
	maxCellineNB         = 256
	maxFailedMeasObjects = 32
)

type IResourceStatusResponseConverter interface {
	Convert(packedBuf []byte) (*models.ResourceStatusResponse, error)
	UnpackX2apPduAsString(packedBuf []byte, maxMessageBufferSize int) (string, error)
}

type ResourceStatusResponseConverter struct {
	X2apPduUnpacker
}

func NewResourceStatusResponseConverter(unpacker X2apPduUnpacker) ResourceStatusResponseConverter {
	return ResourceStatusResponseConverter{unpacker}
}

func buildCellId(cellId C.ECGI_t) string {
	plmnId := C.GoBytes(unsafe.Pointer(cellId.pLMN_Identity.buf), C.int(cellId.pLMN_Identity.size))
	eutranCellIdentifier := C.GoBytes(unsafe.Pointer(cellId.eUTRANcellIdentifier.buf), C.int(cellId.eUTRANcellIdentifier.size))
	return fmt.Sprintf("%x:%x", plmnId, eutranCellIdentifier)
}

func convertMeasurementFailureCauses(measurementFailureCause_List *C.MeasurementFailureCause_List_t, measurementInitiationResult *models.MeasurementInitiationResult) error {
	var MeasurementFailureCauses []*models.MeasurementFailureCause

	count := int(measurementFailureCause_List.list.count)
	if count < 1 || count > maxFailedMeasObjects {
		return fmt.Errorf("invalid number of failure cause elements, %d", int(count))
	}

	measurementFailureCause_ItemIEs_slice := (*[1 << 30]*C.MeasurementFailureCause_ItemIEs_t)(unsafe.Pointer(measurementFailureCause_List.list.array))[:count:count]
	for _, itemIE := range measurementFailureCause_ItemIEs_slice {
		switch itemIE.value.present {
		case C.MeasurementFailureCause_ItemIEs__value_PR_MeasurementFailureCause_Item:
			item := (*C.MeasurementFailureCause_Item_t)(unsafe.Pointer(&itemIE.value.choice[0]))
			measurementFailedReportCharacteristics := C.GoBytes(unsafe.Pointer(item.measurementFailedReportCharacteristics.buf), C.int(item.measurementFailedReportCharacteristics.size))
			measurementFailureCause := models.MeasurementFailureCause{MeasurementFailedReportCharacteristics: measurementFailedReportCharacteristics}
			MeasurementFailureCauses = append(MeasurementFailureCauses, &measurementFailureCause)
			/*Cause ignored - only need to know that the request failed and, possibly, which report characteristics failed*/
		}
		measurementInitiationResult.MeasurementFailureCauses = MeasurementFailureCauses
	}

	return nil
}

func convertMeasurementInitiationResult(measurementInitiationResult_List *C.MeasurementInitiationResult_List_t) ([]*models.MeasurementInitiationResult, error) {
	var measurementInitiationResults []*models.MeasurementInitiationResult

	count := int(measurementInitiationResult_List.list.count)
	if count < 1 || count > maxCellineNB {
		return nil, fmt.Errorf("invalid number of measurement initiation result elements, %d", count)
	}
	measurementInitiationResult_ItemIEs_slice := (*[1 << 30]*C.MeasurementInitiationResult_ItemIEs_t)(unsafe.Pointer(measurementInitiationResult_List.list.array))[:count:count]
	for _, itemIE := range measurementInitiationResult_ItemIEs_slice {
		switch itemIE.value.present {
		case C.MeasurementInitiationResult_ItemIEs__value_PR_MeasurementInitiationResult_Item:
			item := (*C.MeasurementInitiationResult_Item_t)(unsafe.Pointer(&itemIE.value.choice[0]))
			measurementInitiationResult := models.MeasurementInitiationResult{CellId: buildCellId(item.cell_ID)}
			if item.measurementFailureCause_List != nil {
				convertMeasurementFailureCauses(item.measurementFailureCause_List, &measurementInitiationResult)
			}
			measurementInitiationResults = append(measurementInitiationResults, &measurementInitiationResult)
		}

	}

	return measurementInitiationResults, nil
}

func convertResourceStatusIEs(resourceStatusResponse *C.ResourceStatusResponse_t, response *models.ResourceStatusResponse) error {
	count := int(resourceStatusResponse.protocolIEs.list.count)
	resourceStatusResponse_IEs_slice := (*[1 << 30]*C.ResourceStatusResponse_IEs_t)(unsafe.Pointer(resourceStatusResponse.protocolIEs.list.array))[:count:count]
	for _, resourceStatusResponse_IEs := range resourceStatusResponse_IEs_slice {
		switch resourceStatusResponse_IEs.value.present {
		case C.ResourceStatusResponse_IEs__value_PR_Measurement_ID:
			measurement_ID := (*C.Measurement_ID_t)(unsafe.Pointer(&resourceStatusResponse_IEs.value.choice[0]))
			if resourceStatusResponse_IEs.id == C.ProtocolIE_ID_id_ENB1_Measurement_ID {
				response.ENB1_Measurement_ID = int64(*measurement_ID)
			}
			if resourceStatusResponse_IEs.id == C.ProtocolIE_ID_id_ENB2_Measurement_ID {
				response.ENB2_Measurement_ID = int64(*measurement_ID)
			}
		case C.ResourceStatusResponse_IEs__value_PR_CriticalityDiagnostics:
			/*ignored*/

		case C.ResourceStatusResponse_IEs__value_PR_MeasurementInitiationResult_List:
			measurementInitiationResults, err := convertMeasurementInitiationResult((*C.MeasurementInitiationResult_List_t)(unsafe.Pointer(&resourceStatusResponse_IEs.value.choice[0])))
			if err != nil {
				return err
			}
			response.MeasurementInitiationResults = measurementInitiationResults
		}
	}
	return nil
}

func convertResourceStatusResponse(pdu *C.E2AP_PDU_t) (*models.ResourceStatusResponse, error) {
	response := models.ResourceStatusResponse{}

	if pdu.present != C.E2AP_PDU_PR_successfulOutcome {
		return &response, fmt.Errorf("unexpected PDU, %d", int(pdu.present))
	}

	//dereference a union of pointers (C union is represented as a byte array with the size of the largest member)
	successfulOutcome := *(**C.SuccessfulOutcome_t)(unsafe.Pointer(&pdu.choice[0]))
	if successfulOutcome == nil || successfulOutcome.value.present != C.SuccessfulOutcome__value_PR_ResourceStatusResponse {
		return &response, fmt.Errorf("unexpected PDU - not a resource status response")
	}

	resourceStatusResponse := (*C.ResourceStatusResponse_t)(unsafe.Pointer(&successfulOutcome.value.choice[0]))
	if resourceStatusResponse.protocolIEs.list.count == 0 {
		return &response, fmt.Errorf("unexpected PDU - no protocolIEs found")
	}

	if err := convertResourceStatusIEs(resourceStatusResponse, &response); err != nil {
		return &response, err
	}

	return &response, nil
}

// Convert pdu to public ResourceStatusResponse
func (r ResourceStatusResponseConverter) Convert(packedBuf []byte) (*models.ResourceStatusResponse, error) {
	pdu, err := r.UnpackX2apPdu(packedBuf)
	if err != nil {
		return nil, err
	}

	defer C.delete_pdu(pdu)
	return convertResourceStatusResponse(pdu)
}
