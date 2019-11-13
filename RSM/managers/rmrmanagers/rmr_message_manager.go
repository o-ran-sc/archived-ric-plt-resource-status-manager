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
	"rsm/converters"
	"rsm/logger"
	"rsm/models"
	"rsm/providers/rmrmsghandlerprovider"
	"rsm/rmrcgo"
	"rsm/services"
	"rsm/services/rmrsender"
	"time"
)

type RmrMessageManager struct {
	logger          *logger.Logger
	handlerProvider *rmrmsghandlerprovider.MessageHandlerProvider
}

func NewRmrMessageManager(logger *logger.Logger, config *configuration.Configuration, rnibDataService services.RNibDataService, rmrSender *rmrsender.RmrSender, resourceStatusService *services.ResourceStatusService, rsConverter converters.ResourceStatusResponseConverter, rsFailureConverter converters.ResourceStatusFailureConverter) *RmrMessageManager {
	handlerProvider := rmrmsghandlerprovider.NewMessageHandlerProvider(logger, config, rnibDataService, rmrSender, resourceStatusService, rsConverter, rsFailureConverter)

	return &RmrMessageManager{
		handlerProvider: handlerProvider,
		logger:          logger,
	}
}

func (m RmrMessageManager) HandleMessage(mbuf *rmrcgo.MBuf) {

	msgHandler, err := m.handlerProvider.GetMessageHandler(mbuf.MType)

	if err != nil {
		m.logger.Errorf("%s", err)
		return
	}

	request := models.NewRmrRequest(mbuf.Meid, *mbuf.Payload, time.Now())
	go msgHandler.Handle(request)
}
