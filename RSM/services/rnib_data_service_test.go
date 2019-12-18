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


package services

import (
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"net"
	"rsm/configuration"
	"rsm/logger"
	"rsm/mocks"
	"strings"
	"testing"
)

func setupRnibDataServiceTest(t *testing.T) (*rNibDataService, *mocks.RnibReaderMock) {
	return setupRnibDataServiceTestWithMaxAttempts(t, 3)
}

func setupRnibDataServiceTestWithMaxAttempts(t *testing.T, maxAttempts int) (*rNibDataService, *mocks.RnibReaderMock) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}

	config, err := configuration.ParseConfiguration()
	if err != nil {
		t.Errorf("#... - failed to parse configuration error: %s", err)
	}
	config.Rnib.MaxRnibConnectionAttempts = maxAttempts

	rnibReaderMock := &mocks.RnibReaderMock{}
	rsmReaderMock := &mocks.RsmReaderMock{}
	rsmWriterMock := &mocks.RsmWriterMock{}

	rnibDataService := NewRnibDataService(logger, config, rnibReaderMock, rsmReaderMock, rsmWriterMock)
	assert.NotNil(t, rnibDataService)

	return rnibDataService, rnibReaderMock
}

func TestSuccessfulGetNodeb(t *testing.T) {
	rnibDataService, readerMock := setupRnibDataServiceTest(t)

	invName := "abcd"
	nodebInfo := &entities.NodebInfo{}
	readerMock.On("GetNodeb", invName).Return(nodebInfo, nil)

	res, err := rnibDataService.GetNodeb(invName)
	readerMock.AssertNumberOfCalls(t, "GetNodeb", 1)
	assert.Equal(t, nodebInfo, res)
	assert.Nil(t, err)
}

func TestConnFailureGetNodeb(t *testing.T) {
	rnibDataService, readerMock := setupRnibDataServiceTest(t)

	invName := "abcd"
	var nodeb *entities.NodebInfo = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetNodeb", invName).Return(nodeb, mockErr)

	res, err := rnibDataService.GetNodeb(invName)
	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	assert.True(t, strings.Contains(err.Error(), "connection error"))
	assert.Equal(t, nodeb, res)
}

func TestSuccessfulGetNodebIdList(t *testing.T) {
	rnibDataService, readerMock := setupRnibDataServiceTest(t)

	nodeIds := []*entities.NbIdentity{}
	readerMock.On("GetListNodebIds").Return(nodeIds, nil)

	res, err := rnibDataService.GetListNodebIds()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 1)
	assert.Equal(t, nodeIds, res)
	assert.Nil(t, err)
}

func TestConnFailureGetNodebIdList(t *testing.T) {
	rnibDataService, readerMock := setupRnibDataServiceTest(t)

	var nodeIds []*entities.NbIdentity = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res, err := rnibDataService.GetListNodebIds()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 3)
	assert.True(t, strings.Contains(err.Error(), "connection error"))
	assert.Equal(t, nodeIds, res)
}

func TestConnFailureTwiceGetNodebIdList(t *testing.T) {
	rnibDataService, readerMock := setupRnibDataServiceTest(t)

	invName := "abcd"
	var nodeb *entities.NodebInfo = nil
	var nodeIds []*entities.NbIdentity = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetNodeb", invName).Return(nodeb, mockErr)
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res, err := rnibDataService.GetListNodebIds()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 3)
	assert.True(t, strings.Contains(err.Error(), "connection error"))
	assert.Equal(t, nodeIds, res)

	res2, err := rnibDataService.GetNodeb(invName)
	readerMock.AssertNumberOfCalls(t, "GetNodeb", 3)
	assert.True(t, strings.Contains(err.Error(), "connection error"))
	assert.Equal(t, nodeb, res2)
}

func TestConnFailureWithAnotherConfig(t *testing.T) {
	rnibDataService, readerMock := setupRnibDataServiceTestWithMaxAttempts(t, 5)

	var nodeIds []*entities.NbIdentity = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res, err := rnibDataService.GetListNodebIds()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 5)
	assert.True(t, strings.Contains(err.Error(), "connection error"))
	assert.Equal(t, nodeIds, res)
}

func TestPingRnibConnFailure(t *testing.T) {
	rnibDataService, readerMock := setupRnibDataServiceTest(t)

	var nodeIds []*entities.NbIdentity = nil
	mockErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res := rnibDataService.PingRnib()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 3)
	assert.False(t, res)
}

func TestPingRnibOkNoError(t *testing.T) {
	rnibDataService, readerMock := setupRnibDataServiceTest(t)

	var nodeIds []*entities.NbIdentity = nil
	readerMock.On("GetListNodebIds").Return(nodeIds, nil)

	res := rnibDataService.PingRnib()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 1)
	assert.True(t, res)
}

func TestPingRnibOkOtherError(t *testing.T) {
	rnibDataService, readerMock := setupRnibDataServiceTest(t)

	var nodeIds []*entities.NbIdentity = nil
	mockErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	readerMock.On("GetListNodebIds").Return(nodeIds, mockErr)

	res := rnibDataService.PingRnib()
	readerMock.AssertNumberOfCalls(t, "GetListNodebIds", 1)
	assert.True(t, res)
}