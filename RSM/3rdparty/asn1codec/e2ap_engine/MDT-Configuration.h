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

#ifndef	_MDT_Configuration_H_
#define	_MDT_Configuration_H_


#include "asn_application.h"

/* Including external dependencies */
#include "MDT-Activation.h"
#include "AreaScopeOfMDT.h"
#include "MeasurementsToActivate.h"
#include "M1ReportingTrigger.h"
#include "constr_SEQUENCE.h"

#ifdef __cplusplus
extern "C" {
#endif

/* Forward declarations */
struct M1ThresholdEventA2;
struct M1PeriodicReporting;
struct ProtocolExtensionContainer;

/* MDT-Configuration */
typedef struct MDT_Configuration {
	MDT_Activation_t	 mdt_Activation;
	AreaScopeOfMDT_t	 areaScopeOfMDT;
	MeasurementsToActivate_t	 measurementsToActivate;
	M1ReportingTrigger_t	 m1reportingTrigger;
	struct M1ThresholdEventA2	*m1thresholdeventA2;	/* OPTIONAL */
	struct M1PeriodicReporting	*m1periodicReporting;	/* OPTIONAL */
	struct ProtocolExtensionContainer	*iE_Extensions;	/* OPTIONAL */
	/*
	 * This type is extensible,
	 * possible extensions are below.
	 */
	
	/* Context for parsing across buffer boundaries */
	asn_struct_ctx_t _asn_ctx;
} MDT_Configuration_t;

/* Implementation */
extern asn_TYPE_descriptor_t asn_DEF_MDT_Configuration;
extern asn_SEQUENCE_specifics_t asn_SPC_MDT_Configuration_specs_1;
extern asn_TYPE_member_t asn_MBR_MDT_Configuration_1[7];

#ifdef __cplusplus
}
#endif

#endif	/* _MDT_Configuration_H_ */
#include "asn_internal.h"
