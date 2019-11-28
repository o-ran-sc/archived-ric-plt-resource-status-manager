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

#ifndef	_SubbandCQIList_H_
#define	_SubbandCQIList_H_


#include "asn_application.h"

/* Including external dependencies */
#include "asn_SEQUENCE_OF.h"
#include "constr_SEQUENCE_OF.h"

#ifdef __cplusplus
extern "C" {
#endif

/* Forward declarations */
struct SubbandCQIItem;

/* SubbandCQIList */
typedef struct SubbandCQIList {
	A_SEQUENCE_OF(struct SubbandCQIItem) list;
	
	/* Context for parsing across buffer boundaries */
	asn_struct_ctx_t _asn_ctx;
} SubbandCQIList_t;

/* Implementation */
extern asn_TYPE_descriptor_t asn_DEF_SubbandCQIList;
extern asn_SET_OF_specifics_t asn_SPC_SubbandCQIList_specs_1;
extern asn_TYPE_member_t asn_MBR_SubbandCQIList_1[1];
extern asn_per_constraints_t asn_PER_type_SubbandCQIList_constr_1;

#ifdef __cplusplus
}
#endif

#endif	/* _SubbandCQIList_H_ */
#include "asn_internal.h"
