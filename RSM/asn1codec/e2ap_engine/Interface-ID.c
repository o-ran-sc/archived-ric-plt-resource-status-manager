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
 * From ASN.1 module "E2SM-gNB-X2-IEs"
 * 	found in "../../asnFiles/e2sm-gNB-X2-release-1-v041.asn"
 * 	`asn1c -fcompound-names -fincludes-quoted -fno-include-deps -findirect-choice -gen-PER -no-gen-OER -D.`
 */

#include "Interface-ID.h"

#include "GlobalENB-ID.h"
#include "GlobalGNB-ID.h"
asn_per_constraints_t asn_PER_type_Interface_ID_constr_1 CC_NOTUSED = {
	{ APC_CONSTRAINED | APC_EXTENSIBLE,  1,  1,  0,  1 }	/* (0..1,...) */,
	{ APC_UNCONSTRAINED,	-1, -1,  0,  0 },
	0, 0	/* No PER value map */
};
asn_TYPE_member_t asn_MBR_Interface_ID_1[] = {
	{ ATF_POINTER, 0, offsetof(struct Interface_ID, choice.global_eNB_ID),
		(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_GlobalENB_ID,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"global-eNB-ID"
		},
	{ ATF_POINTER, 0, offsetof(struct Interface_ID, choice.global_gNB_ID),
		(ASN_TAG_CLASS_CONTEXT | (1 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_GlobalGNB_ID,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"global-gNB-ID"
		},
};
static const asn_TYPE_tag2member_t asn_MAP_Interface_ID_tag2el_1[] = {
    { (ASN_TAG_CLASS_CONTEXT | (0 << 2)), 0, 0, 0 }, /* global-eNB-ID */
    { (ASN_TAG_CLASS_CONTEXT | (1 << 2)), 1, 0, 0 } /* global-gNB-ID */
};
asn_CHOICE_specifics_t asn_SPC_Interface_ID_specs_1 = {
	sizeof(struct Interface_ID),
	offsetof(struct Interface_ID, _asn_ctx),
	offsetof(struct Interface_ID, present),
	sizeof(((struct Interface_ID *)0)->present),
	asn_MAP_Interface_ID_tag2el_1,
	2,	/* Count of tags in the map */
	0, 0,
	2	/* Extensions start */
};
asn_TYPE_descriptor_t asn_DEF_Interface_ID = {
	"Interface-ID",
	"Interface-ID",
	&asn_OP_CHOICE,
	0,	/* No effective tags (pointer) */
	0,	/* No effective tags (count) */
	0,	/* No tags (pointer) */
	0,	/* No tags (count) */
	{ 0, &asn_PER_type_Interface_ID_constr_1, CHOICE_constraint },
	asn_MBR_Interface_ID_1,
	2,	/* Elements count */
	&asn_SPC_Interface_ID_specs_1	/* Additional specs */
};

