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


    #include <stdbool.h>
    #include <stdio.h>
    #include <stdlib.h>
    #include <resource_status_request_wrapper.h>

    void test_build_pack_resource_status_request();
    void test_unpack(void);

    int
    main(int argc, char* argv[])
    {
    	test_build_pack_resource_status_request();
        exit(0);
    }

    /*reportCharacteristics:
1)    First Bit = PRB Periodic,
2)    Second Bit = TNL load Ind Periodic,
3)    Third Bit = HW Load Ind Periodic,
4)    Fourth Bit = Composite Available Capacity Periodic, this bit should be set to 1 if at least one of the First, Second
or Third bits is set to 1,
5)    Fifth Bit = ABS Status Periodic,
6)    Sixth Bit = RSRP Measurement Report Periodic,
7)    Seventh Bit = CSI Report Periodic
     */
    void test_build_pack_resource_status_request(){
        size_t error_buf_size = 8192;
        size_t packed_buf_size = 4096;
        unsigned char outBuf[packed_buf_size];
        char errorBuf[error_buf_size];
        uint8_t pLMN_Identity[] = {0xa,0xb,0xc};
        uint8_t eUTRANCellIdentifier[] = {0xab, 0xcd, 0x70, 0};
        E2AP_PDU_t *pdu;

        bool result = build_pack_resource_status_request(
        		pLMN_Identity, eUTRANCellIdentifier,
        		15 /*measurement_ID*/, 0  /*measurement_ID2*/,
				Registration_Request_start /*Registration_Request_start,Registration_Request_stop,Registration_Request_partial_stop, Registration_Request_add*/,
        		0xf0,
				ReportingPeriodicity_one_thousand_ms /*ReportingPeriodicity_one_thousand_ms, ReportingPeriodicity_two_thousand_ms, ReportingPeriodicity_five_thousand_ms,ReportingPeriodicity_ten_thousand_ms*/,
				PartialSuccessIndicator_partial_success_allowed /*PartialSuccessIndicator_partial_success_allowed*/,
				ReportingPeriodicityRSRPMR_one_hundred_20_ms /*ReportingPeriodicityRSRPMR_one_hundred_20_ms, ReportingPeriodicityRSRPMR_two_hundred_40_ms, ReportingPeriodicityRSRPMR_four_hundred_80_ms,ReportingPeriodicityRSRPMR_six_hundred_40_ms*/,
				ReportingPeriodicityCSIR_ms5 /*	ReportingPeriodicityCSIR_ms5, ReportingPeriodicityCSIR_ms10,ReportingPeriodicityCSIR_ms20,ReportingPeriodicityCSIR_ms40,ReportingPeriodicityCSIR_ms80*/,
        		&packed_buf_size, outBuf, error_buf_size, errorBuf);

        if (!result) {
            printf("#%s failed. Packing error %s\n", __func__,errorBuf);
            return;
        }
        printf("packed size:%lu\nPayload:\n", packed_buf_size);
        for (size_t i = 0; i < packed_buf_size; ++i)
            printf("%02x",outBuf[i]);
        printf("\n");

        pdu = calloc(1, sizeof(E2AP_PDU_t));
        if (!unpack_pdu_aux(pdu, packed_buf_size, outBuf,error_buf_size, errorBuf,ATS_ALIGNED_BASIC_PER)){
        	printf("#%s failed. Unpacking error %s\n", __func__, errorBuf);
        }
        errorBuf[0] = 0;
        asn1_pdu_printer(pdu, sizeof(errorBuf), errorBuf);
        printf("#%s: %s\n", __func__, errorBuf);

    }

