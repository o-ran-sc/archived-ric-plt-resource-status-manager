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

#include "E-RABs-ToBeModified-SgNBModReqd-Item-SgNBPDCPpresent.h"

#include "E-RAB-Level-QoS-Parameters.h"
#include "ULConfiguration.h"
#include "GTPtunnelEndpoint.h"
#include "ProtocolExtensionContainer.h"
asn_TYPE_member_t asn_MBR_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent_1[] = {
	{ ATF_POINTER, 5, offsetof(struct E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent, requested_MCG_E_RAB_Level_QoS_Parameters),
		(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_E_RAB_Level_QoS_Parameters,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"requested-MCG-E-RAB-Level-QoS-Parameters"
		},
	{ ATF_POINTER, 4, offsetof(struct E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent, uL_Configuration),
		(ASN_TAG_CLASS_CONTEXT | (1 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ULConfiguration,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"uL-Configuration"
		},
	{ ATF_POINTER, 3, offsetof(struct E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent, sgNB_UL_GTP_TEIDatPDCP),
		(ASN_TAG_CLASS_CONTEXT | (2 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_GTPtunnelEndpoint,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"sgNB-UL-GTP-TEIDatPDCP"
		},
	{ ATF_POINTER, 2, offsetof(struct E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent, s1_DL_GTP_TEIDatSgNB),
		(ASN_TAG_CLASS_CONTEXT | (3 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_GTPtunnelEndpoint,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"s1-DL-GTP-TEIDatSgNB"
		},
	{ ATF_POINTER, 1, offsetof(struct E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent, iE_Extensions),
		(ASN_TAG_CLASS_CONTEXT | (4 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ProtocolExtensionContainer_170P73,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"iE-Extensions"
		},
};
static const int asn_MAP_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent_oms_1[] = { 0, 1, 2, 3, 4 };
static const ber_tlv_tag_t asn_DEF_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent_tags_1[] = {
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static const asn_TYPE_tag2member_t asn_MAP_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent_tag2el_1[] = {
    { (ASN_TAG_CLASS_CONTEXT | (0 << 2)), 0, 0, 0 }, /* requested-MCG-E-RAB-Level-QoS-Parameters */
    { (ASN_TAG_CLASS_CONTEXT | (1 << 2)), 1, 0, 0 }, /* uL-Configuration */
    { (ASN_TAG_CLASS_CONTEXT | (2 << 2)), 2, 0, 0 }, /* sgNB-UL-GTP-TEIDatPDCP */
    { (ASN_TAG_CLASS_CONTEXT | (3 << 2)), 3, 0, 0 }, /* s1-DL-GTP-TEIDatSgNB */
    { (ASN_TAG_CLASS_CONTEXT | (4 << 2)), 4, 0, 0 } /* iE-Extensions */
};
asn_SEQUENCE_specifics_t asn_SPC_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent_specs_1 = {
	sizeof(struct E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent),
	offsetof(struct E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent, _asn_ctx),
	asn_MAP_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent_tag2el_1,
	5,	/* Count of tags in the map */
	asn_MAP_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent_oms_1,	/* Optional members */
	5, 0,	/* Root/Additions */
	5,	/* First extension addition */
};
asn_TYPE_descriptor_t asn_DEF_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent = {
	"E-RABs-ToBeModified-SgNBModReqd-Item-SgNBPDCPpresent",
	"E-RABs-ToBeModified-SgNBModReqd-Item-SgNBPDCPpresent",
	&asn_OP_SEQUENCE,
	asn_DEF_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent_tags_1,
	sizeof(asn_DEF_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent_tags_1)
		/sizeof(asn_DEF_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent_tags_1[0]), /* 1 */
	asn_DEF_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent_tags_1,	/* Same as above */
	sizeof(asn_DEF_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent_tags_1)
		/sizeof(asn_DEF_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent_tags_1[0]), /* 1 */
	{ 0, 0, SEQUENCE_constraint },
	asn_MBR_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent_1,
	5,	/* Elements count */
	&asn_SPC_E_RABs_ToBeModified_SgNBModReqd_Item_SgNBPDCPpresent_specs_1	/* Additional specs */
};

