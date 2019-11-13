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

package testhelper

import (
	"rsm/configuration"
	"rsm/logger"
	"rsm/mocks"
	"rsm/rmrcgo"
	"rsm/services"
	"rsm/services/rmrsender"
	"rsm/tests"
	"testing"
)

func InitRmrSender(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *rmrsender.RmrSender {
	rmrMessenger := rmrcgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return rmrsender.NewRmrSender(log, rmrMessenger)
}

func InitLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#tests.initLog - failed to initialize logger, error: %s", err)
	}
	return log
}

func InitTestCase(t *testing.T) (services.RNibDataService, *rmrsender.RmrSender, *logger.Logger) {
	logger := InitLog(t)
	config, err := configuration.ParseConfiguration()
	if err != nil {
		t.Errorf("#tests.InitTestCase - failed to parse configuration, error: %s", err)
	}

	rnibReaderMock := &mocks.RnibReaderMock{}

	rsmReaderMock := &mocks.RsmReaderMock{}
	rsmWriterMock := &mocks.RsmWriterMock{}

	rmrSender := InitRmrSender(&mocks.RmrMessengerMock{}, logger)
	rnibDataService := services.NewRnibDataService(logger, config, rnibReaderMock, rsmReaderMock, rsmWriterMock)
	return rnibDataService, rmrSender, logger
}
