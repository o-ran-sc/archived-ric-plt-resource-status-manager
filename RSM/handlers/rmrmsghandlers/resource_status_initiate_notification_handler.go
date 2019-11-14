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
package rmrmsghandlers

import (
	"encoding/json"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"rsm/enums"
	"rsm/logger"
	"rsm/models"
	"rsm/services"
)

type ResourceStatusInitiateNotificationHandler struct {
	logger                *logger.Logger
	rnibDataService       services.RNibDataService
	resourceStatusService services.IResourceStatusService
	requestName           string
}

func NewResourceStatusInitiateNotificationHandler(logger *logger.Logger, rnibDataService services.RNibDataService, resourceStatusService services.IResourceStatusService, requestName string) ResourceStatusInitiateNotificationHandler {
	return ResourceStatusInitiateNotificationHandler{
		logger:                logger,
		rnibDataService:       rnibDataService,
		resourceStatusService: resourceStatusService,
		requestName:           requestName,
	}
}

func (h ResourceStatusInitiateNotificationHandler) UnmarshalResourceStatusPayload(inventoryName string, payload []byte) (*models.ResourceStatusPayload, error) {
	unmarshalledPayload := models.ResourceStatusPayload{}
	err := json.Unmarshal(payload, &unmarshalledPayload)

	if err != nil {
		h.logger.Errorf("#ResourceStatusInitiateNotificationHandler.UnmarshalResourceStatusPayload - RAN name: %s - Error unmarshaling RMR request payload: %v", inventoryName, err)
		return nil, err
	}

	if unmarshalledPayload.NodeType == entities.Node_UNKNOWN {
		h.logger.Errorf("#ResourceStatusInitiateNotificationHandler.UnmarshalResourceStatusPayload - RAN name: %s - Unknown Node Type", inventoryName)
		return nil, fmt.Errorf("unknown node type for RAN %s", inventoryName)
	}

	h.logger.Infof("#ResourceStatusInitiateNotificationHandler.UnmarshalResourceStatusPayload - Unmarshaled payload successfully: %+v", payload)
	return &unmarshalledPayload, nil

}

func (h ResourceStatusInitiateNotificationHandler) SaveRsmRanInfoStopTrue(inventoryName string) {
	rsmRanInfo := models.NewRsmRanInfo(inventoryName, 0, 0, enums.Stop, true)
	_ = h.rnibDataService.SaveRsmRanInfo(rsmRanInfo)
}

func (h ResourceStatusInitiateNotificationHandler) Handle(request *models.RmrRequest) {
	inventoryName := request.RanName
	h.logger.Infof("#ResourceStatusInitiateNotificationHandler.Handle - RAN name: %s - Received %s notification", inventoryName, h.requestName)

	payload, err := h.UnmarshalResourceStatusPayload(inventoryName, request.Payload)

	if err != nil {
		return
	}

	if payload.NodeType != entities.Node_ENB {
		h.logger.Debugf("#ResourceStatusInitiateNotificationHandler.Handle - RAN name: %s, Node type isn't ENB", inventoryName)
		return
	}

	config, err := h.rnibDataService.GetRsmGeneralConfiguration()

	if err != nil {
		return
	}

	if !config.EnableResourceStatus {
		h.SaveRsmRanInfoStopTrue(inventoryName)
		return
	}

	nodeb, err := h.rnibDataService.GetNodeb(inventoryName)

	if err != nil {
		return
	}

	nodebConnectionStatus := nodeb.GetConnectionStatus()

	if nodebConnectionStatus != entities.ConnectionStatus_CONNECTED {
		h.logger.Errorf("#ResourceStatusInitiateNotificationHandler.Handle - RAN name: %s - RAN's connection status isn't CONNECTED", inventoryName)
		h.SaveRsmRanInfoStopTrue(inventoryName)
		return
	}

	rsmRanInfo := models.NewRsmRanInfo(inventoryName, enums.Enb1MeasurementId, 0, enums.Start, false)
	err = h.rnibDataService.SaveRsmRanInfo(rsmRanInfo)

	if err != nil {
		return
	}

	err = h.resourceStatusService.BuildAndSendInitiateRequest(nodeb, config, rsmRanInfo.Enb1MeasurementId)

	if err != nil {
		h.SaveRsmRanInfoStopTrue(inventoryName)
		return
	}
}
