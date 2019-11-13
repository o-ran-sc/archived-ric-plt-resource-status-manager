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

package mocks

import (
	"github.com/stretchr/testify/mock"
	"rsm/logger"
	"rsm/rmrcgo"
)

type RmrMessengerMock struct {
	mock.Mock
}

func (m *RmrMessengerMock) Init(readyIntervalSec int, port string, maxMsgSize int, flags int, logger *logger.Logger) rmrcgo.RmrMessenger {
	args := m.Called(readyIntervalSec, port, maxMsgSize, flags, logger)
	return args.Get(0).(rmrcgo.RmrMessenger)
}

func (m *RmrMessengerMock) SendMsg(msg *rmrcgo.MBuf) (*rmrcgo.MBuf, error) {
	args := m.Called(msg)
	return args.Get(0).(*rmrcgo.MBuf), args.Error(1)
}

func (m *RmrMessengerMock) RecvMsg() (*rmrcgo.MBuf, error) {
	args := m.Called()
	return args.Get(0).(*rmrcgo.MBuf), args.Error(1)
}

func (m *RmrMessengerMock) RtsMsg(msg *rmrcgo.MBuf) {
	m.Called()
}

func (m *RmrMessengerMock) FreeMsg() {
	m.Called()
}

func (m *RmrMessengerMock) IsReady() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *RmrMessengerMock) Close() {
	m.Called()
}
