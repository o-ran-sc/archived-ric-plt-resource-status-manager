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


package rmrsender

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"rsm/logger"
	"rsm/mocks"
	"rsm/models"
	"rsm/rmrcgo"
	"testing"
)

func initRmrSenderTest(t *testing.T) (*logger.Logger, *mocks.RmrMessengerMock) {
	log := InitLog(t)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrMessengerMock.On("IsReady").Return(true)
	rmrMessengerMock.On("Close").Return()
	return log, rmrMessengerMock
}

func TestRmrSenderSendSuccess(t *testing.T) {
	logger, rmrMessengerMock := initRmrSenderTest(t)

	ranName := "test"
	payload := []byte("some payload")
	xaction := []byte(ranName)
	mbuf := rmrcgo.NewMBuf(123, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock.On("SendMsg", mbuf).Return(&rmrcgo.MBuf{}, nil)
	rmrMsg := models.NewRmrMessage(123, ranName, payload)
	rmrMessenger := rmrcgo.RmrMessenger(rmrMessengerMock)
	rmrSender := NewRmrSender(logger, rmrMessenger)
	err := rmrSender.Send(rmrMsg)
	assert.Nil(t, err)
	rmrMessengerMock.AssertCalled(t, "SendMsg", mbuf)

}

func TestRmrSenderSendFailure(t *testing.T) {
	logger, rmrMessengerMock := initRmrSenderTest(t)

	ranName := "test"
	payload := []byte("some payload")
	xaction := []byte(ranName)
	mbuf := rmrcgo.NewMBuf(123, len(payload), ranName, &payload, &xaction)
	rmrMessengerMock.On("SendMsg", mbuf).Return(mbuf, fmt.Errorf("rmr send failure"))
	rmrMsg := models.NewRmrMessage(123, ranName, payload)
	rmrMessenger := rmrcgo.RmrMessenger(rmrMessengerMock)
	rmrSender := NewRmrSender(logger, rmrMessenger)
	err := rmrSender.Send(rmrMsg)
	rmrMessengerMock.AssertCalled(t, "SendMsg", mbuf)
	assert.NotNil(t, err)
}

func InitLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#tests.initLog - failed to initialize logger, error: %s", err)
	}
	return log
}