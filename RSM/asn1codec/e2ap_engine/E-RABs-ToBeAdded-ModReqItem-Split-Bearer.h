/*
 *
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
 *
 */


/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "X2AP-PDU-Contents"
 * 	found in "../../asnFiles/X2AP-PDU-Contents.asn"
 * 	`asn1c -fcompound-names -fincludes-quoted -fno-include-deps -findirect-choice -gen-PER -no-gen-OER -D.`
 */

#ifndef	_E_RABs_ToBeAdded_ModReqItem_Split_Bearer_H_
#define	_E_RABs_ToBeAdded_ModReqItem_Split_Bearer_H_


#include "asn_application.h"

/* Including external dependencies */
#include "E-RAB-ID.h"
#include "E-RAB-Level-QoS-Parameters.h"
#include "GTPtunnelEndpoint.h"
#include "constr_SEQUENCE.h"

#ifdef __cplusplus
extern "C" {
#endif

/* Forward declarations */
struct ProtocolExtensionContainer;

/* E-RABs-ToBeAdded-ModReqItem-Split-Bearer */
typedef struct E_RABs_ToBeAdded_ModReqItem_Split_Bearer {
	E_RAB_ID_t	 e_RAB_ID;
	E_RAB_Level_QoS_Parameters_t	 e_RAB_Level_QoS_Parameters;
	GTPtunnelEndpoint_t	 meNB_GTPtunnelEndpoint;
	struct ProtocolExtensionContainer	*iE_Extensions;	/* OPTIONAL */
	/*
	 * This type is extensible,
	 * possible extensions are below.
	 */
	
	/* Context for parsing across buffer boundaries */
	asn_struct_ctx_t _asn_ctx;
} E_RABs_ToBeAdded_ModReqItem_Split_Bearer_t;

/* Implementation */
extern asn_TYPE_descriptor_t asn_DEF_E_RABs_ToBeAdded_ModReqItem_Split_Bearer;
extern asn_SEQUENCE_specifics_t asn_SPC_E_RABs_ToBeAdded_ModReqItem_Split_Bearer_specs_1;
extern asn_TYPE_member_t asn_MBR_E_RABs_ToBeAdded_ModReqItem_Split_Bearer_1[4];

#ifdef __cplusplus
}
#endif

#endif	/* _E_RABs_ToBeAdded_ModReqItem_Split_Bearer_H_ */
#include "asn_internal.h"
