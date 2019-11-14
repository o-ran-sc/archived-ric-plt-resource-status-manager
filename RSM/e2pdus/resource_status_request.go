/*******************************************************************************
 *
 *   Copyright (c) 2019 AT&T Intellectual Property.
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 *
 *******************************************************************************/
package e2pdus

// #cgo CFLAGS: -I../asn1codec/inc/ -I../asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../asn1codec/lib/ -L../asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <stdlib.h>
// #include <asn1codec_utils.h>
// #include <resource_status_request_wrapper.h>
import "C"
import (
	"fmt"
	"github.com/pkg/errors"
	"rsm/enums"
	"unsafe"
)

const (
	MaxAsn1PackedBufferSize       = 4096
	MaxAsn1CodecMessageBufferSize = 4096
)

type Measurement_ID int64

type ResourceStatusRequestData struct {
	CellIdList                   []string // PLMNIdentifier:eUTRANCellIdentifier
	MeasurementID                Measurement_ID
	MeasurementID2               Measurement_ID
	PartialSuccessAllowed        bool
	PrbPeriodic                  bool
	TnlLoadIndPeriodic           bool
	HwLoadIndPeriodic            bool
	AbsStatusPeriodic            bool
	RsrpMeasurementPeriodic      bool
	CsiPeriodic                  bool
	PeriodicityMS                enums.ReportingPeriodicity
	PeriodicityRsrpMeasurementMS enums.ReportingPeriodicityRSRPMR
	PeriodicityCsiMS             enums.ReportingPeriodicityCSIR
}

func BuildPackedResourceStatusRequest(registrationRequest enums.Registration_Request, request *ResourceStatusRequestData, maxAsn1PackedBufferSize int, maxAsn1CodecMessageBufferSize int, withDebug bool) ([]byte, string, error) {

	packedBuf := make([]byte, maxAsn1PackedBufferSize)
	errBuf := make([]C.char, maxAsn1CodecMessageBufferSize)
	packedBufSize := C.ulong(len(packedBuf))
	pduAsString := ""

	pLMNIdentities := make([]*C.uchar, len(request.CellIdList))
	eUTRANCellIdentifiers := make([]*C.uchar, len(request.CellIdList))

	for i, cellID := range request.CellIdList {
		var pLMNIdentity, eUTRANCellIdentifier []byte
		if _, err := fmt.Sscanf(cellID, "%x:%x", &pLMNIdentity, &eUTRANCellIdentifier); err != nil {
			return nil, "", fmt.Errorf("BuildPackedResourceStatusRequest() - unexpected CellID value [%s]@%d (want: \"<PLMNIdentifier>:<eUTRANCellIdentifier>\"), err: %s", cellID, i, err)
		}
		pLMNIdentities[i], eUTRANCellIdentifiers[i] = (*C.uchar)(C.CBytes(pLMNIdentity)), (*C.uchar)(C.CBytes(eUTRANCellIdentifier))
	}

	defer func() {
		for _, cPtr := range pLMNIdentities {
			C.free(unsafe.Pointer(cPtr))

		}
		for _, cPtr := range eUTRANCellIdentifiers {
			C.free(unsafe.Pointer(cPtr))
		}

	}()

	/*
		9.2.0	General
		When specifying information elements which are to be represented by bit strings, if not otherwise specifically stated in the semantics description of the concerned IE or elsewhere, the following principle applies with regards to the ordering of bits:
		-	The first bit (leftmost bit) contains the most significant bit (MSB);
		-	The last bit (rightmost bit) contains the least significant bit (LSB);
		-	When importing bit strings from other specifications, the first bit of the bit string contains the first bit of the concerned information.

	*/
	/*reportCharacteristics:
	1)    First Bit = PRB Periodic.
	2)    Second Bit = TNL load Ind Periodic.
	3)    Third Bit = HW Load Ind Periodic.
	4)    Fourth Bit = Composite Available Capacity Periodic, this bit should be set to 1 if at least one of the First, Second or Third bits is set to 1.
	5)    Fifth Bit = ABPartialSuccessIndicator_tS Status Periodic.
	6)    Sixth Bit = RSRP Measurement Report Periodic.
	7)    Seventh Bit = CSI Report Periodic.
	*/
	var prbPeriodic, tnlLoadIndPeriodic, hwLoadIndPeriodic, compositeAvailablCapacityPeriodic, absStatusPeriodic, rsrpMeasurementPeriodic, csiPeriodic int
	var partialSuccessAllowed C.PartialSuccessIndicator_t

	if request.PartialSuccessAllowed {
		partialSuccessAllowed = C.PartialSuccessIndicator_partial_success_allowed
	} else {
		partialSuccessAllowed = C.PartialSuccessIndicator_t(-1)
	}
	if request.PrbPeriodic {
		prbPeriodic = 1
	}
	if request.TnlLoadIndPeriodic {
		tnlLoadIndPeriodic = 1
	}
	if request.HwLoadIndPeriodic {
		hwLoadIndPeriodic = 1
	}
	if request.PrbPeriodic || request.TnlLoadIndPeriodic || request.HwLoadIndPeriodic {
		compositeAvailablCapacityPeriodic = 1
	}
	if request.AbsStatusPeriodic {
		absStatusPeriodic = 1
	}
	if request.RsrpMeasurementPeriodic {
		rsrpMeasurementPeriodic = 1
	}
	if request.CsiPeriodic {
		csiPeriodic = 1
	}
	reportCharacteristics := uint32(prbPeriodic<<7 | tnlLoadIndPeriodic<<6 | hwLoadIndPeriodic<<5 | compositeAvailablCapacityPeriodic<<4 | absStatusPeriodic<<3 | rsrpMeasurementPeriodic<<2 | csiPeriodic<<1)

	if !C.build_pack_resource_status_request(
		(**C.uchar)(unsafe.Pointer(&pLMNIdentities[0])),
		(**C.uchar)(unsafe.Pointer(&eUTRANCellIdentifiers[0])),
		C.ulong(len(request.CellIdList)),
		C.Measurement_ID_t(request.MeasurementID),
		C.Measurement_ID_t(request.MeasurementID2),
		C.Registration_Request_t(registrationRequest),
		C.uint(reportCharacteristics),
		C.ReportingPeriodicity_t(request.PeriodicityMS-1),
		partialSuccessAllowed,
		C.ReportingPeriodicityRSRPMR_t(request.PeriodicityRsrpMeasurementMS-1),
		C.ReportingPeriodicityCSIR_t(request.PeriodicityCsiMS-1),
		&packedBufSize,
		(*C.uchar)(unsafe.Pointer(&packedBuf[0])),
		C.ulong(len(errBuf)),
		&errBuf[0]) {
		return nil, "", errors.New(fmt.Sprintf("BuildPackedResourceStatusRequest - packing error: %s", C.GoString(&errBuf[0])))
	}

	if withDebug {
		pdu := C.new_pdu()
		defer C.delete_pdu(pdu)
		if C.per_unpack_pdu(pdu, packedBufSize, (*C.uchar)(unsafe.Pointer(&packedBuf[0])), C.size_t(len(errBuf)), &errBuf[0]) {
			C.asn1_pdu_printer(pdu, C.size_t(len(errBuf)), &errBuf[0])
			pduAsString = C.GoString(&errBuf[0])
		}
	}

	return packedBuf[:packedBufSize], pduAsString, nil
}
