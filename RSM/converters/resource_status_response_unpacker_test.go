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

func TestResourceStatusResponseConverter(t *testing.T) {
	logger, _ := logger.InitLogger(logger.DebugLevel)
	unpacker := NewX2apPduUnpacker(logger, e2pdus.MaxAsn1CodecMessageBufferSize)
	rsConverters := NewResourceStatusResponseConverter(unpacker)

	var testCases = []struct {
		response  string
		packedPdu string
	}{
		{
			response: "ENB1_Measurement_ID: 15, ENB2_Measurement_ID: 13, MeasurementInitiationResults:[ CellId: 02f829:0007ab50, MeasurementFailureCauses: [  ]  ]",
			/*
				SuccessfulOutcome ::= {
					procedureCode: 9
					criticality: 0 (reject)
					value: ResourceStatusResponse ::= {
						protocolIEs: ProtocolIE-Container ::= {
							ResourceStatusResponse-IEs ::= {
								id: 39
								criticality: 0 (reject)
								value: 15
							}
							ResourceStatusResponse-IEs ::= {
								id: 40
								criticality: 0 (reject)
								value: 13
							}
							ResourceStatusResponse-IEs ::= {
								id: 65
								criticality: 1 (ignore)
								value: MeasurementInitiationResult-List ::= {
									ProtocolIE-Single-Container ::= {
										id: 66
										criticality: 1 (ignore)
										value: MeasurementInitiationResult-Item ::= {
											cell-ID: ECGI ::= {
												pLMN-Identity: 02 F8 29
												eUTRANcellIdentifier: 00 07 AB 50 (4 bits unused)
											}
										}
									}
								}
							}
						}
					}
				}
			*/

			packedPdu: "200900220000030027000300000e0028000300000c0041400d00004240080002f8290007ab50",
		},
		{response: "ENB1_Measurement_ID: 1, ENB2_Measurement_ID: 75, MeasurementInitiationResults:[ CellId: 133023:02b030a0, MeasurementFailureCauses: [ MeasurementFailedReportCharacteristics: 02000000 MeasurementFailedReportCharacteristics: 04000000 MeasurementFailedReportCharacteristics: 08000000 MeasurementFailedReportCharacteristics: 20000000 MeasurementFailedReportCharacteristics: 40000000 MeasurementFailedReportCharacteristics: 80000000  ]  ]",
			/*
				 SuccessfulOutcome ::= {
				    procedureCode: 9
				    criticality: 0 (reject)
				    value: ResourceStatusResponse ::= {
				        protocolIEs: ProtocolIE-Container ::= {
				            ResourceStatusResponse-IEs ::= {
				                id: 39
				                criticality: 0 (reject)
				                value: 1
				            }
				            ResourceStatusResponse-IEs ::= {
				                id: 40
				                criticality: 0 (reject)
				                value: 75
				            }
				            ResourceStatusResponse-IEs ::= {
				                id: 65
				                criticality: 1 (ignore)
				                value: MeasurementInitiationResult-List ::= {
				                    ProtocolIE-Single-Container ::= {
				                        id: 66
				                        criticality: 1 (ignore)
				                        value: MeasurementInitiationResult-Item ::= {
				                            cell-ID: ECGI ::= {
				                                pLMN-Identity: 13 30 23
				                                eUTRANcellIdentifier: 02 B0 30 A0 (4 bits unused)
				                            }
				                            measurementFailureCause-List: MeasurementFailureCause-List ::= {
				                                ProtocolIE-Single-Container ::= {
				                                    id: 67
				                                    criticality: 1 (ignore)
				                                    value: MeasurementFailureCause-Item ::= {
				                                        measurementFailedReportCharacteristics: 02 00 00 00
				                                        cause: 20 (measurement-temporarily-not-available)
				                                    }
				                                }
				                                ProtocolIE-Single-Container ::= {
				                                    id: 67
				                                    criticality: 1 (ignore)
				                                    value: MeasurementFailureCause-Item ::= {
				                                        measurementFailedReportCharacteristics: 04 00 00 00
				                                        cause: 20 (measurement-temporarily-not-available)
				                                    }
				                                }
				                                ProtocolIE-Single-Container ::= {
				                                    id: 67
				                                    criticality: 1 (ignore)
				                                    value: MeasurementFailureCause-Item ::= {
				                                        measurementFailedReportCharacteristics: 08 00 00 00
				                                        cause: 20 (measurement-temporarily-not-available)
				                                    }
				                                }
				                                ProtocolIE-Single-Container ::= {
				                                    id: 67
				                                    criticality: 1 (ignore)
				                                    value: MeasurementFailureCause-Item ::= {
				                                        measurementFailedReportCharacteristics: 20 00 00 00
				                                        cause: 20 (measurement-temporarily-not-available)
				                                    }
				                                }
				                                ProtocolIE-Single-Container ::= {
				                                    id: 67
				                                    criticality: 1 (ignore)
				                                    value: MeasurementFailureCause-Item ::= {
				                                        measurementFailedReportCharacteristics: 40 00 00 00
				                                        cause: 20 (measurement-temporarily-not-available)
				                                    }
				                                }
				                                ProtocolIE-Single-Container ::= {
				                                    id: 67
				                                    criticality: 1 (ignore)
				                                    value: MeasurementFailureCause-Item ::= {
				                                        measurementFailedReportCharacteristics: 80 00 00 00
				                                        cause: 20 (measurement-temporarily-not-available)
				                                    }
				                                }
				                            }
				                        }
				                    }
				                }
				            }
				        }
				    }
				}
			*/
			packedPdu: "20090065000003002700030000000028000300004a00414050000042404b4013302302b030a2800043400700020000000a000043400700040000000a000043400700080000000a000043400700200000000a000043400700400000000a000043400700800000000a00"},
	}

	for _, tc := range testCases {
		t.Run(tc.packedPdu, func(t *testing.T) {
			var payload []byte

			_, err := fmt.Sscanf(tc.packedPdu, "%x", &payload)
			if err != nil {
				t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
			}

			response, err := rsConverters.Convert(payload)
			if err != nil {
				t.Errorf("want: success, got: unpack failed. Error: %v\n", err)
			}

			got := response.String()
			if len(tc.response) != len(got) {
				t.Errorf("\nwant :\t[%s]\n got: \t\t[%s]\n", tc.response, response)
			}
			if strings.Compare(tc.response, got) != 0 {
				t.Errorf("\nwant :\t[%s]\n got: \t\t[%s]\n", tc.response, got)
			}
		})
	}
}

/*unpacking error*/

func TestResourceStatusResponseConverterError(t *testing.T) {
	logger, _ := logger.InitLogger(logger.InfoLevel)
	unpacker := NewX2apPduUnpacker(logger, e2pdus.MaxAsn1CodecMessageBufferSize)
	rsConverters := NewResourceStatusResponseConverter(unpacker)

	wantError := "unpacking error: #src/asn1codec_utils.c.unpack_pdu_aux - Failed to decode E2AP-PDU (consumed 0), error = 0 Success"
	//--------------------2006002a
	inputPayloadAsStr := "2006002b000002001500080002f82900007a8000140017000000630002f8290007ab50102002f829000001000133"
	var payload []byte
	_, err := fmt.Sscanf(inputPayloadAsStr, "%x", &payload)
	if err != nil {
		t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
	}

	_, err = rsConverters.Convert(payload)
	if err != nil {
		if 0 != strings.Compare(fmt.Sprintf("%s", err), wantError) {
			t.Errorf("want failure: %s, got: %s", wantError, err)
		}
	} else {
		t.Errorf("want failure: %s, got: success", wantError)

	}
}

func TestResourceStatusResponseConverterPduOfFailure(t *testing.T) {
	logger, _ := logger.InitLogger(logger.InfoLevel)
	unpacker := NewX2apPduUnpacker(logger, e2pdus.MaxAsn1CodecMessageBufferSize)
	rsConverters := NewResourceStatusResponseConverter(unpacker)

	wantError := "unexpected PDU, 3"
	inputPayloadAsStr := "400900170000030027000300000000280003000049000540020a80"
	var payload []byte
	_, err := fmt.Sscanf(inputPayloadAsStr, "%x", &payload)
	if err != nil {
		t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
	}

	_, err = rsConverters.Convert(payload)
	if err != nil {
		if 0 != strings.Compare(fmt.Sprintf("%s", err), wantError) {
			t.Errorf("want failure: %s, got: %s", wantError, err)
		}
	} else {
		t.Errorf("want failure: %s, got: success", wantError)

	}
}

func TestResourceStatusResponseConverterWrongPdu(t *testing.T) {
	logger, _ := logger.InitLogger(logger.InfoLevel)
	unpacker := NewX2apPduUnpacker(logger,e2pdus.MaxAsn1CodecMessageBufferSize)
	rsConverters := NewResourceStatusResponseConverter(unpacker)

	wantError := "unexpected PDU - not a resource status response"
	inputPayloadAsStr := "2006002a000002001500080002f82900007a8000140017000000630002f8290007ab50102002f829000001000133"
	var payload []byte
	_, err := fmt.Sscanf(inputPayloadAsStr, "%x", &payload)
	if err != nil {
		t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
	}

	_, err = rsConverters.Convert(payload)
	if err != nil {
		if 0 != strings.Compare(fmt.Sprintf("%s", err), wantError) {
			t.Errorf("want failure: %s, got: %s", wantError, err)
		}
	} else {
		t.Errorf("want failure: %s, got: success", wantError)

	}
}
