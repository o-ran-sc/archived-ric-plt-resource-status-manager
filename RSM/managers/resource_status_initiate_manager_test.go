package managers

import (
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"rsm/configuration"
	"rsm/e2pdus"
	"rsm/enums"
	"rsm/logger"
	"rsm/mocks"
	"rsm/rmrcgo"
	"rsm/rsmerrors"
	"rsm/services"
	"rsm/tests/testhelper"
	"testing"
	"time"
)

const RanName = "test"

func initResourceStatusInitiateManagerTest(t *testing.T) (*mocks.RmrMessengerMock, *mocks.RnibReaderMock, *e2pdus.ResourceStatusRequestData, *ResourceStatusInitiateManager) {
	logger, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}

	config, err := configuration.ParseConfiguration()
	if err != nil {
		t.Errorf("#... - failed to parse configuration error: %s", err)
	}

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := testhelper.InitRmrSender(rmrMessengerMock, logger)

	readerMock := &mocks.RnibReaderMock{}

	resourceStatusRequestData := &e2pdus.ResourceStatusRequestData{}
	populateResourceStatusInitiateRequestParams(resourceStatusRequestData, config)

	rnibDataService := services.NewRnibDataService(logger, config, readerMock)
	resourceStatusInitiateManager := NewResourceStatusInitiateManager(logger, rnibDataService, rmrSender)
	return rmrMessengerMock, readerMock, resourceStatusRequestData, resourceStatusInitiateManager
}

func TestGetNodebFailure(t *testing.T) {
	rmrMessengerMock, readerMock, resourceStatusInitiateRequestParams, resourceStatusInitiateManager := initResourceStatusInitiateManagerTest(t)
	var nodebInfo *entities.NodebInfo
	readerMock.On("GetNodeb", RanName).Return(nodebInfo, common.NewInternalError(errors.New("Error")))
	err := resourceStatusInitiateManager.Execute(RanName, resourceStatusInitiateRequestParams)
	readerMock.AssertCalled(t, "GetNodeb", RanName)
	assert.IsType(t, &rsmerrors.RnibDbError{}, err)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestInvalidConnectionStatus(t *testing.T) {
	rmrMessengerMock, readerMock, resourceStatusInitiateRequestParams, resourceStatusInitiateManager := initResourceStatusInitiateManagerTest(t)
	var err error
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}, err)
	err = resourceStatusInitiateManager.Execute(RanName, resourceStatusInitiateRequestParams)
	readerMock.AssertCalled(t, "GetNodeb", RanName)
	assert.IsType(t, &rsmerrors.WrongStateError{}, err)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestPackFailure(t *testing.T) {
	rmrMessengerMock, readerMock, resourceRequestData, resourceStatusInitiateManager := initResourceStatusInitiateManagerTest(t)
	var err error
	nodebInfo := &entities.NodebInfo{
		RanName:          RanName,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Configuration: &entities.NodebInfo_Enb{
			Enb: &entities.Enb{
				ServedCells: []*entities.ServedCellInfo{{CellId: ""}},
			},
		},
	}

	readerMock.On("GetNodeb", RanName).Return(nodebInfo, err)
	err = resourceStatusInitiateManager.Execute(RanName, resourceRequestData)
	assert.Nil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", RanName)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestOneCellSuccess(t *testing.T) {
	cellId := "02f829:0007ab00"
	rmrMessengerMock, readerMock, resourceRequestData, resourceStatusInitiateManager := initResourceStatusInitiateManagerTest(t)
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

	readerMock.On("GetNodeb", RanName).Return(nodebInfo, err)
	expectedPayload := getPackedPayloadForCell(cellId, 1, *resourceRequestData)
	expectedMbuf := rmrcgo.NewMBuf(rmrcgo.RicResStatusReq, len(expectedPayload), RanName, &expectedPayload, &xaction)
	rmrMessengerMock.On("SendMsg", expectedMbuf).Return(&rmrcgo.MBuf{}, err)
	err = resourceStatusInitiateManager.Execute(RanName, resourceRequestData)
	time.Sleep(100 * time.Millisecond)
	readerMock.AssertCalled(t, "GetNodeb", RanName)
	assert.Nil(t, err)
	rmrMessengerMock.AssertCalled(t, "SendMsg", expectedMbuf)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestTwoCellOneFailureOneSuccess(t *testing.T) {
	cellId1 := "02f829:0007ab00"
	cellId2 := "02f829:0007ab50"
	rmrMessengerMock, readerMock, resourceRequestData, resourceStatusInitiateManager := initResourceStatusInitiateManagerTest(t)
	xaction := []byte(RanName)
	var err error
	nodebInfo := &entities.NodebInfo{
		RanName:          RanName,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Configuration: &entities.NodebInfo_Enb{
			Enb: &entities.Enb{
				ServedCells: []*entities.ServedCellInfo{{CellId: cellId1}, {CellId: cellId2}},
			},
		},
	}

	readerMock.On("GetNodeb", RanName).Return(nodebInfo, err)
	expectedPayload1 := getPackedPayloadForCell(cellId1, 1, *resourceRequestData)
	expectedMbuf1 := rmrcgo.NewMBuf(rmrcgo.RicResStatusReq, len(expectedPayload1), RanName, &expectedPayload1, &xaction)

	expectedPayload2 := getPackedPayloadForCell(cellId2, 2, *resourceRequestData)
	expectedMbuf2 := rmrcgo.NewMBuf(rmrcgo.RicResStatusReq, len(expectedPayload2), RanName, &expectedPayload2, &xaction)
	rmrMessengerMock.On("SendMsg", expectedMbuf1).Return(&rmrcgo.MBuf{}, rsmerrors.NewRmrError())
	rmrMessengerMock.On("SendMsg", expectedMbuf2).Return(&rmrcgo.MBuf{}, err)
	err = resourceStatusInitiateManager.Execute(RanName, resourceRequestData)
	time.Sleep(100 * time.Millisecond)
	readerMock.AssertCalled(t, "GetNodeb", RanName)
	assert.Nil(t, err)
	rmrMessengerMock.AssertCalled(t, "SendMsg", expectedMbuf1)
	rmrMessengerMock.AssertCalled(t, "SendMsg", expectedMbuf2)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 2)
}

func TestFiveCellsSuccess(t *testing.T) {
	cellId1 := "02f829:0007ab00"
	cellId2 := "02f829:0007ab50"
	cellId3 := "02f829:0007ab60"
	cellId4 := "02f829:0007ab70"
	cellId5 := "02f829:0007ab80"

	rmrMessengerMock, readerMock, resourceRequestData, resourceStatusInitiateManager := initResourceStatusInitiateManagerTest(t)
	xaction := []byte(RanName)
	var err error
	nodebInfo := &entities.NodebInfo{
		RanName:          RanName,
		ConnectionStatus: entities.ConnectionStatus_CONNECTED,
		Configuration: &entities.NodebInfo_Enb{
			Enb: &entities.Enb{
				ServedCells: []*entities.ServedCellInfo{{CellId: cellId1}, {CellId: cellId2}, {CellId: cellId3}, {CellId: cellId4}, {CellId: cellId5}},
			},
		},
	}

	readerMock.On("GetNodeb", RanName).Return(nodebInfo, err)
	expectedPayload1 := getPackedPayloadForCell(cellId1, 1, *resourceRequestData)
	expectedMbuf1 := rmrcgo.NewMBuf(rmrcgo.RicResStatusReq, len(expectedPayload1), RanName, &expectedPayload1, &xaction)

	expectedPayload2 := getPackedPayloadForCell(cellId2, 2, *resourceRequestData)
	expectedMbuf2 := rmrcgo.NewMBuf(rmrcgo.RicResStatusReq, len(expectedPayload2), RanName, &expectedPayload2, &xaction)

	expectedPayload3 := getPackedPayloadForCell(cellId3, 3, *resourceRequestData)
	expectedMbuf3 := rmrcgo.NewMBuf(rmrcgo.RicResStatusReq, len(expectedPayload3), RanName, &expectedPayload3, &xaction)

	expectedPayload4 := getPackedPayloadForCell(cellId4, 4, *resourceRequestData)
	expectedMbuf4 := rmrcgo.NewMBuf(rmrcgo.RicResStatusReq, len(expectedPayload4), RanName, &expectedPayload4, &xaction)

	expectedPayload5 := getPackedPayloadForCell(cellId5, 5, *resourceRequestData)
	expectedMbuf5 := rmrcgo.NewMBuf(rmrcgo.RicResStatusReq, len(expectedPayload5), RanName, &expectedPayload5, &xaction)

	rmrMessengerMock.On("SendMsg", expectedMbuf1).Return(&rmrcgo.MBuf{}, err)
	rmrMessengerMock.On("SendMsg", expectedMbuf2).Return(&rmrcgo.MBuf{}, err)
	rmrMessengerMock.On("SendMsg", expectedMbuf3).Return(&rmrcgo.MBuf{}, err)
	rmrMessengerMock.On("SendMsg", expectedMbuf4).Return(&rmrcgo.MBuf{}, err)
	rmrMessengerMock.On("SendMsg", expectedMbuf5).Return(&rmrcgo.MBuf{}, err)

	err = resourceStatusInitiateManager.Execute(RanName, resourceRequestData)
	time.Sleep(100 * time.Millisecond)
	readerMock.AssertCalled(t, "GetNodeb", RanName)
	assert.Nil(t, err)
	rmrMessengerMock.AssertCalled(t, "SendMsg", expectedMbuf1)
	rmrMessengerMock.AssertCalled(t, "SendMsg", expectedMbuf2)
	rmrMessengerMock.AssertCalled(t, "SendMsg", expectedMbuf3)
	rmrMessengerMock.AssertCalled(t, "SendMsg", expectedMbuf4)
	rmrMessengerMock.AssertCalled(t, "SendMsg", expectedMbuf5)

	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 5)
}

func getPackedPayloadForCell(cellId string, index int, resourceStatusRequestData e2pdus.ResourceStatusRequestData) []byte {
	resourceStatusRequestData.CellID = cellId
	resourceStatusRequestData.MeasurementID = e2pdus.Measurement_ID(index)
	expectedPayload, _, _ := e2pdus.BuildPackedResourceStatusRequest(enums.Registration_Request_start, &resourceStatusRequestData, e2pdus.MaxAsn1PackedBufferSize, e2pdus.MaxAsn1CodecMessageBufferSize, false)
	return expectedPayload
}

func populateResourceStatusInitiateRequestParams(params *e2pdus.ResourceStatusRequestData, config *configuration.Configuration) {
	params.PartialSuccessAllowed = config.ResourceStatusParams.PartialSuccessAllowed
	params.PrbPeriodic = config.ResourceStatusParams.PrbPeriodic
	params.TnlLoadIndPeriodic = config.ResourceStatusParams.TnlLoadIndPeriodic
	params.HwLoadIndPeriodic = config.ResourceStatusParams.HwLoadIndPeriodic
	params.AbsStatusPeriodic = config.ResourceStatusParams.AbsStatusPeriodic
	params.RsrpMeasurementPeriodic = config.ResourceStatusParams.RsrpMeasurementPeriodic
	params.CsiPeriodic = config.ResourceStatusParams.CsiPeriodic
	params.PeriodicityMS = config.ResourceStatusParams.PeriodicityMs
	params.PeriodicityRsrpMeasurementMS = config.ResourceStatusParams.PeriodicityRsrpMeasurementMs
	params.PeriodicityCsiMS = config.ResourceStatusParams.PeriodicityCsiMs
}
