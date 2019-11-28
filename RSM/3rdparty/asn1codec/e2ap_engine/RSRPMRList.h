/*
 * Copyright 2019 AT&T Intellectual Property
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
 */

/*
 * This source code is part of the near-RT RIC (RAN Intelligent Controller)
 * platform project (RICP).
 */



/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "X2AP-IEs"
 * 	found in "../../asnFiles/X2AP-IEs.asn"
 * 	`asn1c -fcompound-names -fincludes-quoted -fno-include-deps -findirect-choice -gen-PER -no-gen-OER -D.`
 */

#ifndef	_RSRPMRList_H_
#define	_RSRPMRList_H_


#include "asn_application.h"

/* Including external dependencies */
#include "asn_SEQUENCE_OF.h"
#include "RSRPMeasurementResult.h"
#include "constr_SEQUENCE.h"
#include "constr_SEQUENCE_OF.h"

#ifdef __cplusplus
extern "C" {
#endif

/* Forward declarations */
struct ProtocolExtensionContainer;

/* Forward definitions */
typedef struct RSRPMRList__Member {
	RSRPMeasurementResult_t	 rSRPMeasurementResult;
	struct ProtocolExtensionContainer	*iE_Extensions;	/* OPTIONAL */
	/*
	 * This type is extensible,
	 * possible extensions are below.
	 */
	
	/* Context for parsing across buffer boundaries */
	asn_struct_ctx_t _asn_ctx;
} RSRPMRList__Member;

/* RSRPMRList */
typedef struct RSRPMRList {
	A_SEQUENCE_OF(RSRPMRList__Member) list;
	
	/* Context for parsing across buffer boundaries */
	asn_struct_ctx_t _asn_ctx;
} RSRPMRList_t;

/* Implementation */
extern asn_TYPE_descriptor_t asn_DEF_RSRPMRList;
extern asn_SET_OF_specifics_t asn_SPC_RSRPMRList_specs_1;
extern asn_TYPE_member_t asn_MBR_RSRPMRList_1[1];
extern asn_per_constraints_t asn_PER_type_RSRPMRList_constr_1;

#ifdef __cplusplus
}
#endif

#endif	/* _RSRPMRList_H_ */
#include "asn_internal.h"
