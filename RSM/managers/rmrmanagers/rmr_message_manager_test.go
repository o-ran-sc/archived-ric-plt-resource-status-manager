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
//
package rmrmanagers

import (
	"rsm/configuration"
	"rsm/managers"
	"rsm/rmrcgo"
	"rsm/tests/testhelper"
	"testing"
)

func TestRmrMessageManagerSuccess(t *testing.T) {

	rnibDataService, rmrSender, logger := testhelper.InitTestCase(t)
	config, _ := configuration.ParseConfiguration()
	resourceStatusInitiateManager := managers.NewResourceStatusInitiateManager(logger, rnibDataService, rmrSender)

	manager := NewRmrMessageManager(logger, config, rnibDataService, rmrSender, resourceStatusInitiateManager, nil)

	xactionByteArr := []byte("1111111")
	payloadByteArr := []byte("payload")
	msg := rmrcgo.NewMBuf(rmrcgo.RanConnected, len(payloadByteArr), "test", &payloadByteArr, &xactionByteArr)

	manager.HandleMessage(msg)
}

func TestRmrMessageManagerFailure(t *testing.T) {

	rnibDataService, rmrSender, logger := testhelper.InitTestCase(t)
	config, _ := configuration.ParseConfiguration()
	resourceStatusInitiateManager := managers.NewResourceStatusInitiateManager(logger, rnibDataService, rmrSender)

	manager := NewRmrMessageManager(logger, config, rnibDataService, rmrSender, resourceStatusInitiateManager, nil)

	xactionByteArr := []byte("1111111")
	payloadByteArr := []byte("payload")
	msg := rmrcgo.NewMBuf(11, len(payloadByteArr), "test", &payloadByteArr, &xactionByteArr)

	manager.HandleMessage(msg)
}
