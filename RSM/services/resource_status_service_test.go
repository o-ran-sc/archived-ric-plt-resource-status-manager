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

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package services

import (
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"rsm/enums"
	"rsm/logger"
	"rsm/mocks"
	"rsm/models"
	"rsm/rmrcgo"
	"rsm/rsmerrors"
	"rsm/services/rmrsender"
	"rsm/tests"
	"testing"
)

const RanName = "test"
const NodebOneCellPackedExample = "0009003c00000800270003000000001c00010000260004fe000000001d400d00001f40080002f8290007ab00001e4001000040400100006d4001400091400120"
const NodebTwoCellsPackedExample = "0009004800000800270003000000001c00010000260004fe000000001d401901001f40080002f8290007ab00001f40080002f8290007ab50001e4001000040400100006d4001400091400120"
const StopPackedExample = "0009004f0000090027000300000000280003000001001c00014000260004fe000000001d401901001f40080002f8290007ab00001f40080002f8290007ab50001e4001000040400100006d4001400091400120"

func initResourceStatusServiceTest(t *testing.T) (*mocks.RmrMessengerMock, *models.RsmGeneralConfiguration, *ResourceStatusService) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := InitRmrSender(rmrMessengerMock, logger)
	resourceStatusService := NewResourceStatusService(logger, rmrSender)

	rsmGeneralConfiguration := models.RsmGeneralConfiguration{
		EnableResourceStatus:         true,
		PartialSuccessAllowed:        true,
		PrbPeriodic:                  true,
		TnlLoadIndPeriodic:           true,
		HwLoadIndPeriodic:            true,
		AbsStatusPeriodic:            true,
		RsrpMeasurementPeriodic:      true,
		CsiPeriodic:                  true,
		PeriodicityMs:                enums.ReportingPeriodicity_one_thousand_ms,
		PeriodicityRsrpMeasurementMs: enums.ReportingPeriodicityRSRPMR_four_hundred_80_ms,
		PeriodicityCsiMs:             enums.ReportingPeriodicityCSIR_ms20,
	}

	return rmrMessengerMock, &rsmGeneralConfiguration, resourceStatusService
}

func TestOneCellSuccess(t *testing.T) {
	cellId := "02f829:0007ab00"
	rmrMessengerMock, rsmGeneralConfiguration, resourceStatusService := initResourceStatusServiceTest(t)

	xaction := []byte(RanName)
	nodebInfo := &entities.NodebInfo{
		RanName:          RanName,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Configuration: &entities.NodebInfo_Enb{
			Enb: &entities.Enb{
				ServedCells: []*entities.ServedCellInfo{{CellId: cellId}},
			},
		},
	}

	var expectedPayload []byte
	_, _ = fmt.Sscanf(NodebOneCellPackedExample, "%x", &expectedPayload)
	var err error
	expectedMbuf := rmrcgo.NewMBuf(rmrcgo.RicResStatusReq, len(expectedPayload), RanName, &expectedPayload, &xaction)
	rmrMessengerMock.On("SendMsg", expectedMbuf).Return(&rmrcgo.MBuf{}, err)
	err = resourceStatusService.BuildAndSendInitiateRequest(nodebInfo, rsmGeneralConfiguration, enums.Enb1MeasurementId)
	assert.Nil(t, err)
	rmrMessengerMock.AssertCalled(t, "SendMsg", expectedMbuf)
}

func TestTwoCellsSuccess(t *testing.T) {
	cellId1 := "02f829:0007ab00"
	cellId2 := "02f829:0007ab50"
	rmrMessengerMock, rsmGeneralConfiguration, resourceStatusService := initResourceStatusServiceTest(t)
	xaction := []byte(RanName)
	nodebInfo := &entities.NodebInfo{
		RanName:          RanName,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Configuration: &entities.NodebInfo_Enb{
			Enb: &entities.Enb{
				ServedCells: []*entities.ServedCellInfo{{CellId: cellId1}, {CellId: cellId2}},
			},
		},
	}

	var expectedPayload []byte
	_, _ = fmt.Sscanf(NodebTwoCellsPackedExample, "%x", &expectedPayload)
	expectedMbuf := rmrcgo.NewMBuf(rmrcgo.RicResStatusReq, len(expectedPayload), RanName, &expectedPayload, &xaction)
	var err error
	rmrMessengerMock.On("SendMsg", expectedMbuf).Return(&rmrcgo.MBuf{}, err)
	err = resourceStatusService.BuildAndSendInitiateRequest(nodebInfo, rsmGeneralConfiguration, enums.Enb1MeasurementId)
	assert.Nil(t, err)
	rmrMessengerMock.AssertCalled(t, "SendMsg", expectedMbuf)
}

func TestOneCellSendFailure(t *testing.T) {
	cellId := "02f829:0007ab00"
	rmrMessengerMock, rsmGeneralConfiguration, resourceStatusService := initResourceStatusServiceTest(t)

	xaction := []byte(RanName)
	var err error
	nodebInfo := &entities.NodebInfo{
		RanName:          RanName,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Configuration: &entities.NodebInfo_Enb{
			Enb: &entities.Enb{
				ServedCells: []*entities.ServedCellInfo{{CellId: cellId}},
			},
		},
	}

	var expectedPayload []byte
	_, _ = fmt.Sscanf(NodebOneCellPackedExample, "%x", &expectedPayload)
	expectedMbuf := rmrcgo.NewMBuf(rmrcgo.RicResStatusReq, len(expectedPayload), RanName, &expectedPayload, &xaction)
	rmrMessengerMock.On("SendMsg", expectedMbuf).Return(&rmrcgo.MBuf{}, rsmerrors.NewRmrError())
	err = resourceStatusService.BuildAndSendInitiateRequest(nodebInfo, rsmGeneralConfiguration, 1)
	assert.NotNil(t, err)
	rmrMessengerMock.AssertCalled(t, "SendMsg", expectedMbuf)
}



func TestNodebConfigurationFailure(t *testing.T) {
	rmrMessengerMock, rsmGeneralConfiguration, resourceStatusService := initResourceStatusServiceTest(t)
	nodebInfo := &entities.NodebInfo{
		RanName:          RanName,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
	}

	err := resourceStatusService.BuildAndSendInitiateRequest(nodebInfo, rsmGeneralConfiguration, enums.Enb1MeasurementId)
	assert.NotNil(t, err)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestNodebEmptyCellList(t *testing.T) {
	rmrMessengerMock, rsmGeneralConfiguration, resourceStatusService := initResourceStatusServiceTest(t)
	nodebInfo := &entities.NodebInfo{
		RanName:          RanName,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Configuration: &entities.NodebInfo_Enb{
			Enb: &entities.Enb{
				ServedCells: []*entities.ServedCellInfo{},
			},
		},
	}

	err := resourceStatusService.BuildAndSendInitiateRequest(nodebInfo, rsmGeneralConfiguration, enums.Enb1MeasurementId)
	assert.NotNil(t, err)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestPackFailure(t *testing.T) {
	rmrMessengerMock, rsmGeneralConfiguration, resourceStatusService := initResourceStatusServiceTest(t)
	nodebInfo := &entities.NodebInfo{
		RanName:          RanName,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Configuration: &entities.NodebInfo_Enb{
			Enb: &entities.Enb{
				ServedCells: []*entities.ServedCellInfo{{CellId: ""}},
			},
		},
	}

	err := resourceStatusService.BuildAndSendInitiateRequest(nodebInfo, rsmGeneralConfiguration, enums.Enb1MeasurementId)
	assert.NotNil(t, err)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestBuildAndSendStopRequestSuccess(t *testing.T) {
	rmrMessengerMock, rsmGeneralConfiguration, resourceStatusService := initResourceStatusServiceTest(t)

	cellId1 := "02f829:0007ab00"
	cellId2 := "02f829:0007ab50"
	nodebInfo := &entities.NodebInfo{
		RanName:          RanName,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Configuration: &entities.NodebInfo_Enb{
			Enb: &entities.Enb{
				ServedCells: []*entities.ServedCellInfo{{CellId: cellId1}, {CellId: cellId2}},
			},
		},
	}
	xaction := []byte(RanName)
	var expectedPayload []byte
	_, _ = fmt.Sscanf(StopPackedExample, "%x", &expectedPayload)
	expectedMbuf := rmrcgo.NewMBuf(rmrcgo.RicResStatusReq, len(expectedPayload), RanName, &expectedPayload, &xaction)
	var err error
	rmrMessengerMock.On("SendMsg", expectedMbuf).Return(&rmrcgo.MBuf{}, err)
	err = resourceStatusService.BuildAndSendStopRequest(nodebInfo, rsmGeneralConfiguration, enums.Enb1MeasurementId, 2)
	assert.Nil(t, err)
	rmrMessengerMock.AssertCalled(t, "SendMsg", expectedMbuf)
}

func TestBuildAndSendStopRequestSendFailure(t *testing.T) {
	rmrMessengerMock, rsmGeneralConfiguration, resourceStatusService := initResourceStatusServiceTest(t)

	xaction := []byte(RanName)
	cellId1 := "02f829:0007ab00"
	cellId2 := "02f829:0007ab50"
	nodebInfo := &entities.NodebInfo{
		RanName:          RanName,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Configuration: &entities.NodebInfo_Enb{
			Enb: &entities.Enb{
				ServedCells: []*entities.ServedCellInfo{{CellId: cellId1}, {CellId: cellId2}},
			},
		},
	}

	var err error
	var expectedPayload []byte
	_, _ = fmt.Sscanf(StopPackedExample, "%x", &expectedPayload)
	expectedMbuf := rmrcgo.NewMBuf(rmrcgo.RicResStatusReq, len(expectedPayload), RanName, &expectedPayload, &xaction)
	rmrMessengerMock.On("SendMsg", expectedMbuf).Return(&rmrcgo.MBuf{}, rsmerrors.NewRmrError())
	err = resourceStatusService.BuildAndSendStopRequest(nodebInfo, rsmGeneralConfiguration, enums.Enb1MeasurementId, 2)

	assert.NotNil(t, err)
	rmrMessengerMock.AssertCalled(t, "SendMsg", expectedMbuf)
}

func InitRmrSender(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *rmrsender.RmrSender {
	rmrMessenger := rmrcgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return rmrsender.NewRmrSender(log, rmrMessenger)
}
