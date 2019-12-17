/*
 *
 * Copyright 2019 AT&T Intellectual Property
 * Copyright 2019 Nokia
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

 /*
  * This source code is part of the near-RT RIC (RAN Intelligent Controller)
  * platform project (RICP).
  */

#include <string.h>
#include <errno.h>
#undef NDEBUG
#include <assert.h>
#include <InitiatingMessage.h>
#include <ProtocolIE-ID.h>
#include <resource_status_request_wrapper.h>

#define pLMN_Identity_size 3
#define eUTRANcellIdentifier_size 28

static void assignPLMN_Identity (PLMN_Identity_t *pLMN_Identity, uint8_t const* pLMNId);
static void assignEUTRANcellIdentifier (EUTRANCellIdentifier_t *eUTRANCellIdentifier, uint8_t const* eUTRANCellId);
/*
 * Build and pack resource status request.
 * Abort the process on allocation failure.
 *  packed_buf_size - in: size of packed_buf; out: number of chars used.
 */

bool
build_pack_resource_status_request(
		uint8_t const** pLMN_Identities,
		uint8_t const** eUTRANCellIdentifiers,
		size_t nECGIs,
		Measurement_ID_t measurement_ID, Measurement_ID_t measurement_ID2,
		Registration_Request_t registration_Request /*Registration_Request_start,Registration_Request_stop,Registration_Request_partial_stop, Registration_Request_add*/,
		uint32_t reportCharacteristics,
		ReportingPeriodicity_t reportingPeriodicity /*ReportingPeriodicity_one_thousand_ms,	ReportingPeriodicity_two_thousand_ms, ReportingPeriodicity_five_thousand_ms,ReportingPeriodicity_ten_thousand_ms*/,
		PartialSuccessIndicator_t partialSuccessIndicator /*PartialSuccessIndicator_partial_success_allowed*/,
		ReportingPeriodicityRSRPMR_t reportingPeriodicityRSRPMR /*ReportingPeriodicityRSRPMR_one_hundred_20_ms,	ReportingPeriodicityRSRPMR_two_hundred_40_ms, ReportingPeriodicityRSRPMR_four_hundred_80_ms,ReportingPeriodicityRSRPMR_six_hundred_40_ms*/,
		ReportingPeriodicityCSIR_t reportingPeriodicityCSIR /*ReportingPeriodicityCSIR_ms5, ReportingPeriodicityCSIR_ms10,ReportingPeriodicityCSIR_ms20,ReportingPeriodicityCSIR_ms40,ReportingPeriodicityCSIR_ms80*/,
		size_t* packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf
)
{
	return build_pack_resource_status_request_aux(
			pLMN_Identities,
            eUTRANCellIdentifiers,
            nECGIs,
			measurement_ID,	measurement_ID2,
			registration_Request,
			reportCharacteristics,
			reportingPeriodicity,
			partialSuccessIndicator,
			reportingPeriodicityRSRPMR,
			reportingPeriodicityCSIR,
			packed_buf_size, packed_buf,err_buf_size,err_buf,ATS_ALIGNED_BASIC_PER);

}

bool
build_pack_resource_status_request_aux(
		uint8_t const** pLMN_Identities,
		uint8_t const** eUTRANCellIdentifiers,
		size_t nECGIs,
		Measurement_ID_t measurement_ID, Measurement_ID_t measurement_ID2,
		Registration_Request_t registration_Request /*Registration_Request_start,Registration_Request_stop,Registration_Request_partial_stop, Registration_Request_add*/,
		uint32_t reportCharacteristics,
		ReportingPeriodicity_t reportingPeriodicity /*ReportingPeriodicity_one_thousand_ms,	ReportingPeriodicity_two_thousand_ms, ReportingPeriodicity_five_thousand_ms,ReportingPeriodicity_ten_thousand_ms*/,
		PartialSuccessIndicator_t partialSuccessIndicator /*PartialSuccessIndicator_partial_success_allowed*/,
		ReportingPeriodicityRSRPMR_t reportingPeriodicityRSRPMR /*ReportingPeriodicityRSRPMR_one_hundred_20_ms,	ReportingPeriodicityRSRPMR_two_hundred_40_ms, ReportingPeriodicityRSRPMR_four_hundred_80_ms,ReportingPeriodicityRSRPMR_six_hundred_40_ms*/,
		ReportingPeriodicityCSIR_t reportingPeriodicityCSIR /*	ReportingPeriodicityCSIR_ms5, ReportingPeriodicityCSIR_ms10,ReportingPeriodicityCSIR_ms20,ReportingPeriodicityCSIR_ms40,ReportingPeriodicityCSIR_ms80*/,
		size_t* packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf,enum asn_transfer_syntax syntax
)
{
	bool rc = true;
	E2AP_PDU_t *pdu = calloc(1, sizeof(E2AP_PDU_t));
	InitiatingMessage_t *initiatingMessage = calloc(1, sizeof(InitiatingMessage_t));
	ResourceStatusRequest_t *request;

    assert(pdu != 0);
    assert(initiatingMessage != 0);


    pdu->present = E2AP_PDU_PR_initiatingMessage;
    pdu->choice.initiatingMessage = initiatingMessage;

    initiatingMessage->procedureCode = ProcedureCode_id_resourceStatusReportingInitiation;
    initiatingMessage->criticality = Criticality_reject;
    initiatingMessage->value.present = InitiatingMessage__value_PR_ResourceStatusRequest;
    request = &initiatingMessage->value.choice.ResourceStatusRequest;

    ResourceStatusRequest_IEs_t *measurementID_ie = calloc(1, sizeof(ResourceStatusRequest_IEs_t));
    assert(measurementID_ie != 0);
    ASN_SEQUENCE_ADD(&request->protocolIEs, measurementID_ie);

    measurementID_ie->id = ProtocolIE_ID_id_ENB1_Measurement_ID;
    measurementID_ie->criticality = Criticality_reject;
    measurementID_ie->value.present = ResourceStatusRequest_IEs__value_PR_Measurement_ID;
    measurementID_ie->value.choice.Measurement_ID = measurement_ID;


    if (registration_Request == Registration_Request_stop
    		|| registration_Request == Registration_Request_partial_stop
			|| registration_Request == Registration_Request_add) {
		ResourceStatusRequest_IEs_t *measurementID2_ie = calloc(1, sizeof(ResourceStatusRequest_IEs_t));
		assert(measurementID2_ie != 0);
		ASN_SEQUENCE_ADD(&request->protocolIEs, measurementID2_ie);

		measurementID2_ie->id = ProtocolIE_ID_id_ENB2_Measurement_ID;
		measurementID2_ie->criticality = Criticality_reject;
		measurementID2_ie->value.present = ResourceStatusRequest_IEs__value_PR_Measurement_ID;
		measurementID2_ie->value.choice.Measurement_ID = measurement_ID2;
    }

    ResourceStatusRequest_IEs_t *registration_Request_ie = calloc(1, sizeof(ResourceStatusRequest_IEs_t));
    assert(registration_Request_ie != 0);
    ASN_SEQUENCE_ADD(&request->protocolIEs, registration_Request_ie);

    registration_Request_ie->id = ProtocolIE_ID_id_Registration_Request;
    registration_Request_ie->criticality = Criticality_reject;
    registration_Request_ie->value.present = ResourceStatusRequest_IEs__value_PR_Registration_Request;
    registration_Request_ie->value.choice.Registration_Request = registration_Request;

    if (reportCharacteristics){
		ResourceStatusRequest_IEs_t *reportCharacteristics_ie = calloc(1, sizeof(ResourceStatusRequest_IEs_t));
		assert(reportCharacteristics_ie != 0);
		ASN_SEQUENCE_ADD(&request->protocolIEs, reportCharacteristics_ie);

		reportCharacteristics_ie->id = ProtocolIE_ID_id_ReportCharacteristics;
		reportCharacteristics_ie->criticality = Criticality_reject;
		reportCharacteristics_ie->value.present = ResourceStatusRequest_IEs__value_PR_ReportCharacteristics;
		reportCharacteristics_ie->value.choice.ReportCharacteristics.size = sizeof(uint32_t);
		reportCharacteristics_ie->value.choice.ReportCharacteristics.bits_unused = 0;
		reportCharacteristics_ie->value.choice.ReportCharacteristics.buf = calloc(1, reportCharacteristics_ie->value.choice.ReportCharacteristics.size);

		assert(reportCharacteristics_ie->value.choice.ReportCharacteristics.buf != 0);
		memcpy(reportCharacteristics_ie->value.choice.ReportCharacteristics.buf, &reportCharacteristics, reportCharacteristics_ie->value.choice.ReportCharacteristics.size);
    }

    ResourceStatusRequest_IEs_t *cellToReport_List_ie = calloc(1, sizeof(ResourceStatusRequest_IEs_t));
    assert(cellToReport_List_ie != 0);
    ASN_SEQUENCE_ADD(&request->protocolIEs, cellToReport_List_ie);

    cellToReport_List_ie->id = ProtocolIE_ID_id_CellToReport;
    cellToReport_List_ie->criticality = Criticality_ignore;
    cellToReport_List_ie->value.present = ResourceStatusRequest_IEs__value_PR_CellToReport_List;

    for (size_t i = 0; i < nECGIs; i++){
        CellToReport_ItemIEs_t *item = calloc(1, sizeof(CellToReport_ItemIEs_t));
        assert(item != 0);
        ASN_SEQUENCE_ADD(&cellToReport_List_ie->value.choice.CellToReport_List, item);

        item->id = ProtocolIE_ID_id_CellToReport_Item;
        item->criticality = Criticality_ignore;
        item->value.present = CellToReport_ItemIEs__value_PR_CellToReport_Item;
        assignPLMN_Identity(&item->value.choice.CellToReport_Item.cell_ID.pLMN_Identity, pLMN_Identities[i]);
        assignEUTRANcellIdentifier(&item->value.choice.CellToReport_Item.cell_ID.eUTRANcellIdentifier,eUTRANCellIdentifiers[i]);
    }

	if (reportingPeriodicity >= 0){
		ResourceStatusRequest_IEs_t *reportingPeriodicity_ie = calloc(1, sizeof(ResourceStatusRequest_IEs_t));
		assert(reportingPeriodicity_ie != 0);
		ASN_SEQUENCE_ADD(&request->protocolIEs, reportingPeriodicity_ie);

		reportingPeriodicity_ie->id = ProtocolIE_ID_id_ReportingPeriodicity;
		reportingPeriodicity_ie->criticality = Criticality_ignore;
		reportingPeriodicity_ie->value.present = ResourceStatusRequest_IEs__value_PR_ReportingPeriodicity;
		reportingPeriodicity_ie->value.choice.ReportingPeriodicity = reportingPeriodicity;
	}

	if (partialSuccessIndicator >= 0){
		ResourceStatusRequest_IEs_t *partialSuccessIndicator_ie = calloc(1, sizeof(ResourceStatusRequest_IEs_t));
		assert(partialSuccessIndicator_ie != 0);
		ASN_SEQUENCE_ADD(&request->protocolIEs, partialSuccessIndicator_ie);

		partialSuccessIndicator_ie->id = ProtocolIE_ID_id_PartialSuccessIndicator;
		partialSuccessIndicator_ie->criticality = Criticality_ignore;
		partialSuccessIndicator_ie->value.present = ResourceStatusRequest_IEs__value_PR_PartialSuccessIndicator;
		partialSuccessIndicator_ie->value.choice.PartialSuccessIndicator = partialSuccessIndicator;
    }

    if (reportingPeriodicityRSRPMR >= 0){
		ResourceStatusRequest_IEs_t *reportingPeriodicityRSRPMR_ie = calloc(1, sizeof(ResourceStatusRequest_IEs_t));
		assert(reportingPeriodicityRSRPMR_ie != 0);
		ASN_SEQUENCE_ADD(&request->protocolIEs, reportingPeriodicityRSRPMR_ie);

		reportingPeriodicityRSRPMR_ie->id = ProtocolIE_ID_id_ReportingPeriodicityRSRPMR;
		reportingPeriodicityRSRPMR_ie->criticality = Criticality_ignore;
		reportingPeriodicityRSRPMR_ie->value.present = ResourceStatusRequest_IEs__value_PR_ReportingPeriodicityRSRPMR;
		reportingPeriodicityRSRPMR_ie->value.choice.ReportingPeriodicityRSRPMR = reportingPeriodicityRSRPMR;
    }

    if (reportingPeriodicityCSIR >= 0){
		ResourceStatusRequest_IEs_t *reportingPeriodicityCSIR_ie = calloc(1, sizeof(ResourceStatusRequest_IEs_t));
		assert(reportingPeriodicityCSIR_ie != 0);
		ASN_SEQUENCE_ADD(&request->protocolIEs, reportingPeriodicityCSIR_ie);

		reportingPeriodicityCSIR_ie->id = ProtocolIE_ID_id_ReportingPeriodicityCSIR;
		reportingPeriodicityCSIR_ie->criticality = Criticality_ignore;
		reportingPeriodicityCSIR_ie->value.present = ResourceStatusRequest_IEs__value_PR_ReportingPeriodicityCSIR;
		reportingPeriodicityCSIR_ie->value.choice.ReportingPeriodicityCSIR = reportingPeriodicityCSIR;
    }

    rc = pack_pdu_aux(pdu, packed_buf_size, packed_buf,err_buf_size, err_buf,syntax);

	ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, pdu);
	return rc;
}

static void assignPLMN_Identity (PLMN_Identity_t *pLMN_Identity, uint8_t const* pLMNId)
{
	pLMN_Identity->size = pLMN_Identity_size;
	pLMN_Identity->buf = calloc(1,pLMN_Identity->size);
	assert(pLMN_Identity->buf != 0);
	memcpy(pLMN_Identity->buf, pLMNId, pLMN_Identity->size);
}

// Assume that eUTRANCellId value is already pushed to the left
static void assignEUTRANcellIdentifier (EUTRANCellIdentifier_t *eUTRANCellIdentifier, uint8_t const* eUTRANCellId)
{
	size_t size_in_bytes = (eUTRANcellIdentifier_size / 8) + ((eUTRANcellIdentifier_size % 8) > 0);
	int unused_bits = 8 - (eUTRANcellIdentifier_size % 8);

	eUTRANCellIdentifier->size = size_in_bytes;
	eUTRANCellIdentifier->bits_unused = unused_bits;
	eUTRANCellIdentifier->buf = calloc(1, eUTRANCellIdentifier->size);
	assert(eUTRANCellIdentifier->buf != 0);
	memcpy(eUTRANCellIdentifier->buf, eUTRANCellId, size_in_bytes);
}
