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
package converters

import (
	"fmt"
	"rsm/e2pdus"
	"rsm/logger"
	"strings"
	"testing"
)

/*
 * Unpack a response returned from RAN.
 * Verify it matches the want pdu.
 */

func TestUnpackX2apPduResponse(t *testing.T) {
	logger, _ := logger.InitLogger(logger.DebugLevel)
	unpacker := NewX2apPduUnpacker(logger, e2pdus.MaxAsn1CodecMessageBufferSize)

	wantPduAsStr := `UnsuccessfulOutcome ::= {
    procedureCode: 9
    criticality: 0 (reject)
    value: ResourceStatusFailure ::= {
        protocolIEs: ProtocolIE-Container ::= {
            ResourceStatusFailure-IEs ::= {
                id: 39
                criticality: 0 (reject)
                value: 15
            }
            ResourceStatusFailure-IEs ::= {
                id: 40
                criticality: 0 (reject)
                value: 13
            }
            ResourceStatusFailure-IEs ::= {
                id: 5
                criticality: 1 (ignore)
                value: 1 (hardware-failure)
            }
            ResourceStatusFailure-IEs ::= {
                id: 68
                criticality: 1 (ignore)
                value: CompleteFailureCauseInformation-List ::= {
                    ProtocolIE-Single-Container ::= {
                        id: 69
                        criticality: 1 (ignore)
                        value: CompleteFailureCauseInformation-Item ::= {
                            cell-ID: ECGI ::= {
                                pLMN-Identity: 02 F8 29
                                eUTRANcellIdentifier: 00 07 AB 50 (4 bits unused)
                            }
                            measurementFailureCause-List: MeasurementFailureCause-List ::= {
                                ProtocolIE-Single-Container ::= {
                                    id: 67
                                    criticality: 1 (ignore)
                                    value: MeasurementFailureCause-Item ::= {
                                        measurementFailedReportCharacteristics: 00 00 00 07
                                        cause: 0 (transfer-syntax-error)
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}`

	inputPayloadAsStr := "400900320000040027000300000e0028000300000c00054001620044401800004540130002f8290007ab500000434006000000000740"
	var payload []byte

	_, err := fmt.Sscanf(inputPayloadAsStr, "%x", &payload)
	if err != nil {
		t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
	}

	response, err := unpacker.UnpackX2apPduAsString(payload, e2pdus.MaxAsn1CodecMessageBufferSize /*message buffer*/)
	if err != nil {
		t.Errorf("want: success, got: unpack failed. Error: %v\n", err)
	}

	want := strings.Fields(wantPduAsStr)
	got := strings.Fields(response)
	if len(want) != len(got) {
		t.Errorf("\nwant :\t[%s]\n got: \t\t[%s]\n", wantPduAsStr, response)
	}
	for i := 0; i < len(want); i++ {
		if strings.Compare(want[i], got[i]) != 0 {
			t.Errorf("\nwant :\t[%s]\n got: \t\t[%s]\n", wantPduAsStr, strings.TrimSpace(response))
		}

	}
}

/*unpacking error*/

func TestUnpackX2apPduError(t *testing.T) {
	logger, _ := logger.InitLogger(logger.InfoLevel)
	unpacker := NewX2apPduUnpacker(logger, e2pdus.MaxAsn1CodecMessageBufferSize)

	wantError := "unpacking error: #src/asn1codec_utils.c.unpack_pdu_aux - Failed to decode E2AP-PDU (consumed 0), error = 0 Success"
	//--------------------2006002a
	inputPayloadAsStr := "2006002b000002001500080002f82900007a8000140017000000630002f8290007ab50102002f829000001000133"
	var payload []byte
	_, err := fmt.Sscanf(inputPayloadAsStr, "%x", &payload)
	if err != nil {
		t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
	}

	_, err = unpacker.UnpackX2apPduAsString(payload, e2pdus.MaxAsn1CodecMessageBufferSize /*message buffer*/)
	if err != nil {
		if 0 != strings.Compare(fmt.Sprintf("%s", err), wantError) {
			t.Errorf("want failure: %s, got: %s", wantError, err)
		}
	} else {
		t.Errorf("want failure: %s, got: success", wantError)

	}
}
