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

package services

import (
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"rsm/e2pdus"
	"rsm/enums"
	"rsm/logger"
	"rsm/models"
	"rsm/rmrcgo"
	"rsm/services/rmrsender"
)

type ResourceStatusService struct {
	logger    *logger.Logger
	rmrSender *rmrsender.RmrSender
}

type IResourceStatusService interface {
	BuildAndSendInitiateRequest(nodeb *entities.NodebInfo, config *models.RsmGeneralConfiguration, enb1MeasurementId int64) error
	//BuildAndSendStopRequest(config *models.RsmGeneralConfiguration, ranName string, enb1MeasurementId int64, enb2MeasurementId int64) error
}

func NewResourceStatusService(logger *logger.Logger, rmrSender *rmrsender.RmrSender) *ResourceStatusService {
	return &ResourceStatusService{
		logger:    logger,
		rmrSender: rmrSender,
	}
}

func (m *ResourceStatusService) BuildAndSendInitiateRequest(nodeb *entities.NodebInfo, config *models.RsmGeneralConfiguration, enb1MeasurementId int64) error {

	cellIdList, err := m.extractCellIdList(nodeb)

	if err != nil {
		return err
	}

	requestParams := buildResourceStatusInitiateRequestParams(config, cellIdList, enb1MeasurementId)

	payload, payloadAsString, err := e2pdus.BuildPackedResourceStatusRequest(enums.Registration_Request_start, requestParams, e2pdus.MaxAsn1PackedBufferSize, e2pdus.MaxAsn1CodecMessageBufferSize, m.logger.DebugEnabled())

	if err != nil {
		m.logger.Errorf("#ResourceStatusService.BuildAndSendInitiateRequest - RAN name: %s. Failed to build and pack resource status initiate request. error: %s", nodeb.RanName, err)
		return err
	}

	m.logger.Debugf("#ResourceStatusService.BuildAndSendInitiateRequest - RAN name: %s. Successfully build packed payload: %s", nodeb.RanName, payloadAsString)
	rmrMsg := models.NewRmrMessage(rmrcgo.RicResStatusReq, nodeb.RanName, payload)

	return m.rmrSender.Send(rmrMsg)
}

/*func (m *ResourceStatusService) BuildAndSendStopRequest(config *models.RsmGeneralConfiguration, ranName string, enb1MeasurementId int64, enb2MeasurementId int64) error {

	requestParams := &e2pdus.ResourceStatusRequestData{
		MeasurementID:  e2pdus.Measurement_ID(enb1MeasurementId),
		MeasurementID2: e2pdus.Measurement_ID(enb2MeasurementId),
	}

	payload, payloadAsString, err := e2pdus.BuildPackedResourceStatusRequest(enums.Registration_Request_stop, requestParams, e2pdus.MaxAsn1PackedBufferSize, e2pdus.MaxAsn1CodecMessageBufferSize, m.logger.DebugEnabled())

	if err != nil {
		m.logger.Errorf("#ResourceStatusService.BuildAndSendStopRequest - RAN name: %s. Failed to build and pack resource status stop request. error: %s", ranName, err)
		return err
	}

	m.logger.Debugf("#ResourceStatusService.BuildAndSendStopRequest - RAN name: %s. Successfully build packed payload: %s", ranName, payloadAsString)
	rmrMsg := models.NewRmrMessage(rmrcgo.RicResStatusReq, ranName, payload)

	return m.rmrSender.Send(rmrMsg)
}*/

func (m *ResourceStatusService) extractCellIdList(nodeb *entities.NodebInfo) ([]string, error) {

	enb, ok := nodeb.Configuration.(*entities.NodebInfo_Enb)

	if !ok {
		m.logger.Errorf("#ResourceStatusService.extractCellIdList - RAN name: %s - invalid configuration", nodeb.RanName)
		return []string{}, fmt.Errorf("Invalid configuration for RAN %s", nodeb.RanName)
	}

	cells := enb.Enb.ServedCells

	if len(cells) == 0 {
		m.logger.Errorf("#ResourceStatusService.extractCellIdList - RAN name: %s - empty cell list", nodeb.RanName)
		return []string{}, fmt.Errorf("empty cell list for RAN %s", nodeb.RanName)
	}

	cellIdList := make([]string, len(cells))
	for index, cellInfo := range cells {
		cellIdList[index] = cellInfo.CellId
	}

	return cellIdList, nil
}

func buildResourceStatusInitiateRequestParams(config *models.RsmGeneralConfiguration, cellIdList []string, enb1MeasurementId int64) *e2pdus.ResourceStatusRequestData {
	return &e2pdus.ResourceStatusRequestData{
		CellIdList:                   cellIdList,
		MeasurementID:                e2pdus.Measurement_ID(enb1MeasurementId),
		PartialSuccessAllowed:        config.PartialSuccessAllowed,
		PrbPeriodic:                  config.PrbPeriodic,
		TnlLoadIndPeriodic:           config.TnlLoadIndPeriodic,
		HwLoadIndPeriodic:            config.HwLoadIndPeriodic,
		AbsStatusPeriodic:            config.AbsStatusPeriodic,
		RsrpMeasurementPeriodic:      config.RsrpMeasurementPeriodic,
		CsiPeriodic:                  config.CsiPeriodic,
		PeriodicityMS:                config.PeriodicityMs,
		PeriodicityRsrpMeasurementMS: config.PeriodicityRsrpMeasurementMs,
		PeriodicityCsiMS:             config.PeriodicityCsiMs,
	}
}
