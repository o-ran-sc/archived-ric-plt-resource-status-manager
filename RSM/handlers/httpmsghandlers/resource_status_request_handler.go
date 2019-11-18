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
package httpmsghandlers

import (
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"rsm/enums"
	"rsm/logger"
	"rsm/models"
	"rsm/rsmerrors"
	"rsm/services"
)

type ResourceStatusRequestHandler struct {
	rNibDataService       services.RNibDataService
	logger                *logger.Logger
	resourceStatusService services.IResourceStatusService
}

func NewResourceStatusRequestHandler(logger *logger.Logger, rNibDataService services.RNibDataService, resourceStatusService services.IResourceStatusService) *ResourceStatusRequestHandler {
	return &ResourceStatusRequestHandler{
		resourceStatusService: resourceStatusService,
		rNibDataService:       rNibDataService,
		logger:                logger,
	}
}

func (h ResourceStatusRequestHandler) Handle(request models.Request) error {

	resourceStatusRequest := request.(models.ResourceStatusRequest)
	config, err := h.rNibDataService.GetRsmGeneralConfiguration()
	if err != nil {
		return rsmerrors.NewRnibDbError()
	}

	config.EnableResourceStatus = resourceStatusRequest.EnableResourceStatus
	err = h.rNibDataService.SaveRsmGeneralConfiguration(config)
	if err != nil {
		return rsmerrors.NewRnibDbError()
	}

	nbIdentityList, err := h.rNibDataService.GetListEnbIds()
	if err != nil {
		return rsmerrors.NewRnibDbError()
	}

	numberOfFails := 0
	for _, nbIdentity := range nbIdentityList {

		nodeb, err := h.rNibDataService.GetNodeb((*nbIdentity).GetInventoryName())
		if err != nil {
			numberOfFails++
			continue
		}

		if nodeb.ConnectionStatus != entities.ConnectionStatus_CONNECTED {
			h.logger.Infof("#ResourceStatusRequestHandler.Handle - RAN name: %s - connection status not CONNECTED, ignore", nodeb.RanName)
			numberOfFails++
			continue
		}

		err = h.saveAndSendRsmRanInfo(nodeb, config)
		if err != nil {
			numberOfFails++
		}
	}

	if numberOfFails > 0 {
		return rsmerrors.NewRsmError(numberOfFails)
	}
	return nil
}

func (h ResourceStatusRequestHandler) saveAndSendRsmRanInfo(nodebInfo *entities.NodebInfo, config *models.RsmGeneralConfiguration) error {

	rsmRanInfo, err := h.rNibDataService.GetRsmRanInfo(nodebInfo.RanName)
	if err != nil {
		return err
	}

	if config.EnableResourceStatus {
		err := h.handleNotStartedRsmRanInfo(nodebInfo, rsmRanInfo, config)
		return err
	}

	//err = h.handleNotStoppedRsmRanInfo(nodebInfo, rsmRanInfo, config)
	return nil
}

/*func (h ResourceStatusRequestHandler) handleNotStoppedRsmRanInfo(nodebInfo *entities.NodebInfo, rsmRanInfo *models.RsmRanInfo, config *models.RsmGeneralConfiguration) error {
	if rsmRanInfo.Action == enums.Stop && rsmRanInfo.ActionStatus {
		return nil
	}

	if rsmRanInfo.Action != enums.Stop {

		err := h.saveRsmRanInfoStopFalse(rsmRanInfo)
		if err != nil {
			return err
		}
	}

	err := h.resourceStatusService.BuildAndSendStopRequest(config, rsmRanInfo.RanName, rsmRanInfo.Enb1MeasurementId, rsmRanInfo.Enb2MeasurementId)
	return err
}*/

func (h ResourceStatusRequestHandler) handleNotStartedRsmRanInfo(nodebInfo *entities.NodebInfo, rsmRanInfo *models.RsmRanInfo, config *models.RsmGeneralConfiguration) error {
	if rsmRanInfo.Action == enums.Start && rsmRanInfo.ActionStatus {
		return nil
	}

	if rsmRanInfo.Action != enums.Start {

		err := h.saveRsmRanInfoStartFalse(rsmRanInfo)
		if err != nil {
			return err
		}
	}

	err := h.resourceStatusService.BuildAndSendInitiateRequest(nodebInfo, config, rsmRanInfo.Enb1MeasurementId)
	return err
}

/*func (h ResourceStatusRequestHandler) saveRsmRanInfoStopFalse(rsmRanInfo *models.RsmRanInfo) error {
	rsmRanInfo.Action = enums.Stop
	rsmRanInfo.ActionStatus = false

	err := h.rNibDataService.SaveRsmRanInfo(rsmRanInfo)
	if err != nil {
		h.logger.Errorf("#ResourceStatusRequestHandler.saveRsmRanInfoStopFalse - failed to save rsm ran data to RNIB. Error: %s", err.Error())
		return err
	}
	return nil
}*/

func (h ResourceStatusRequestHandler) saveRsmRanInfoStartFalse(rsmRanInfo *models.RsmRanInfo) error {
	rsmRanInfo.Action = enums.Start
	rsmRanInfo.ActionStatus = false
	rsmRanInfo.Enb2MeasurementId = 0

	err := h.rNibDataService.SaveRsmRanInfo(rsmRanInfo)
	if err != nil {
		return err
	}
	return nil
}
