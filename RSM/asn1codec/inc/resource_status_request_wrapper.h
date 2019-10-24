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


#include <stddef.h>
#include <stdbool.h>
#include <stdint.h>
#include <asn1codec_utils.h>

#ifndef INC_RESOURCE_STATUS_REQUEST_WRAPPER_H
#define INC_RESOURCE_STATUS_REQUEST_WRAPPER_H

#ifdef __cplusplus
extern "C"
{
#endif

bool
build_pack_resource_status_request(
		uint8_t const* pLMN_Identity,
        uint8_t const* eUTRANCellIdentifier,
		Measurement_ID_t measurement_ID, Measurement_ID_t measurement_ID2,
		Registration_Request_t registration_Request /*Registration_Request_start,Registration_Request_stop,Registration_Request_partial_stop, Registration_Request_add*/,
		uint32_t reportCharacteristics,
		ReportingPeriodicity_t reportingPeriodicity /*ReportingPeriodicity_one_thousand_ms,	ReportingPeriodicity_two_thousand_ms, ReportingPeriodicity_five_thousand_ms,ReportingPeriodicity_ten_thousand_ms*/,
		PartialSuccessIndicator_t partialSuccessIndicator /*PartialSuccessIndicator_partial_success_allowed*/,
		ReportingPeriodicityRSRPMR_t reportingPeriodicityRSRPMR /*ReportingPeriodicityRSRPMR_one_hundred_20_ms,	ReportingPeriodicityRSRPMR_two_hundred_40_ms, ReportingPeriodicityRSRPMR_four_hundred_80_ms,ReportingPeriodicityRSRPMR_six_hundred_40_ms*/,
		ReportingPeriodicityCSIR_t reportingPeriodicityCSIR /*	ReportingPeriodicityCSIR_ms5, ReportingPeriodicityCSIR_ms10,ReportingPeriodicityCSIR_ms20,ReportingPeriodicityCSIR_ms40,ReportingPeriodicityCSIR_ms80*/,
		size_t* packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf
);

bool
build_pack_resource_status_request_aux(
		uint8_t const* pLMN_Identity,
        uint8_t const* eUTRANCellIdentifier,
		Measurement_ID_t measurement_ID, Measurement_ID_t measurement_ID2,
		Registration_Request_t registration_Request /*Registration_Request_start,Registration_Request_stop,Registration_Request_partial_stop, Registration_Request_add*/,
		uint32_t reportCharacteristics,
		ReportingPeriodicity_t reportingPeriodicity /*ReportingPeriodicity_one_thousand_ms,	ReportingPeriodicity_two_thousand_ms, ReportingPeriodicity_five_thousand_ms,ReportingPeriodicity_ten_thousand_ms*/,
		PartialSuccessIndicator_t partialSuccessIndicator /*PartialSuccessIndicator_partial_success_allowed*/,
		ReportingPeriodicityRSRPMR_t reportingPeriodicityRSRPMR /*ReportingPeriodicityRSRPMR_one_hundred_20_ms,	ReportingPeriodicityRSRPMR_two_hundred_40_ms, ReportingPeriodicityRSRPMR_four_hundred_80_ms,ReportingPeriodicityRSRPMR_six_hundred_40_ms*/,
		ReportingPeriodicityCSIR_t reportingPeriodicityCSIR /*	ReportingPeriodicityCSIR_ms5, ReportingPeriodicityCSIR_ms10,ReportingPeriodicityCSIR_ms20,ReportingPeriodicityCSIR_ms40,ReportingPeriodicityCSIR_ms80*/,
		size_t* packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf,enum asn_transfer_syntax syntax
);

#ifdef __cplusplus
}
#endif

#endif /* INC_RESOURCE_STATUS_REQUEST_WRAPPER_H */
 
