//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).


package rmrmsghandlerprovider

import (
	"fmt"
	"rsm/configuration"
	"rsm/converters"
	"rsm/e2pdus"
	"rsm/handlers/rmrmsghandlers"
	"rsm/services"
	"rsm/tests/testhelper"
	"strings"
	"testing"

	"rsm/rmrcgo"
)

/*
 * Verify support for known providers.
 */
func TestGetNotificationHandlerSuccess(t *testing.T) {

	rnibDataService, rmrSender, logger := testhelper.InitTestCase(t)
	config, err := configuration.ParseConfiguration()
	if err != nil {
		t.Errorf("#... - failed to parse configuration error: %s", err)
	}
	resourceStatusService := services.NewResourceStatusService(logger, rmrSender)
	unpacker := converters.NewX2apPduUnpacker(logger, e2pdus.MaxAsn1CodecMessageBufferSize)
	var testCases = []struct {
		msgType int
		handler rmrmsghandlers.RmrMessageHandler
	}{
		{rmrcgo.RanConnected, rmrmsghandlers.NewResourceStatusInitiateNotificationHandler(logger, rnibDataService, resourceStatusService, RanConnected)},
		{rmrcgo.RanRestarted, rmrmsghandlers.NewResourceStatusInitiateNotificationHandler(logger, rnibDataService, resourceStatusService, RanRestarted)},
	}

	for _, tc := range testCases {
		provider := NewMessageHandlerProvider(logger, config, rnibDataService, rmrSender, resourceStatusService, converters.NewResourceStatusResponseConverter(unpacker), converters.NewResourceStatusFailureConverter(unpacker))
		t.Run(fmt.Sprintf("%d", tc.msgType), func(t *testing.T) {
			handler, err := provider.GetMessageHandler(tc.msgType)
			if err != nil {
				t.Errorf("want: handler for message type %d, got: error %s", tc.msgType, err)
			}
			if strings.Compare(fmt.Sprintf("%T", handler), fmt.Sprintf("%T", tc.handler)) != 0 {
				t.Errorf("want: handler %T for message type %d, got: %T", tc.handler, tc.msgType, handler)
			}
		})
	}
}

/*
 * Verify handling of a request for an unsupported message.
 */

func TestGetNotificationHandlerFailure(t *testing.T) {

	rnibDataService, rmrSender, logger := testhelper.InitTestCase(t)
	config, err := configuration.ParseConfiguration()
	if err != nil {
		t.Errorf("#... - failed to parse configuration error: %s", err)
	}
	resourceStatusService := services.NewResourceStatusService(logger, rmrSender)
	unpacker := converters.NewX2apPduUnpacker(logger, e2pdus.MaxAsn1CodecMessageBufferSize)

	var testCases = []struct {
		msgType   int
		errorText string
	}{
		{9999 /*unknown*/, "notification handler not found"},
	}
	for _, tc := range testCases {
		provider := NewMessageHandlerProvider(logger, config, rnibDataService, rmrSender, resourceStatusService, converters.NewResourceStatusResponseConverter(unpacker), converters.NewResourceStatusFailureConverter(unpacker))
		t.Run(fmt.Sprintf("%d", tc.msgType), func(t *testing.T) {
			_, err := provider.GetMessageHandler(tc.msgType)
			if err == nil {
				t.Errorf("want: no handler for message type %d, got: success", tc.msgType)
			}
			if !strings.Contains(fmt.Sprintf("%s", err), tc.errorText) {
				t.Errorf("want: error [%s] for message type %d, got: %s", tc.errorText, tc.msgType, err)
			}
		})
	}
}
