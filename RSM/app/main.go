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


package main

import (
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"gerrit.o-ran-sc.org/r/ric-plt/sdlgo"
	"os"
	"rsm/configuration"
	"rsm/controllers"
	"rsm/converters"
	"rsm/e2pdus"
	"rsm/httpserver"
	"rsm/logger"
	"rsm/managers/rmrmanagers"
	"rsm/providers/httpmsghandlerprovider"
	"rsm/rmrcgo"
	"rsm/rsmdb"
	"rsm/services"
	"rsm/services/rmrreceiver"
	"rsm/services/rmrsender"
	"strconv"
)

func main() {
	config, err := configuration.ParseConfiguration()
	if err != nil {
		fmt.Printf("#app.main - failed to parse configuration, error: %s", err)
		os.Exit(1)
	}
	logLevel, _ := logger.LogLevelTokenToLevel(config.Logging.LogLevel)
	logger, err := logger.InitLogger(logLevel)
	if err != nil {
		fmt.Printf("#app.main - failed to initialize logger, error: %s", err)
		os.Exit(1)
	}
	db := sdlgo.NewDatabase()
	e2mSdl := sdlgo.NewSdlInstance("e2Manager", db)
	rsmSdl := sdlgo.NewSdlInstance("rsm", db)

	defer e2mSdl.Close()
	defer rsmSdl.Close()

	logger.Infof("#app.main - Configuration %s", config)
	rnibDataService := services.NewRnibDataService(logger, config, reader.GetRNibReader(e2mSdl), rsmdb.GetRsmReader(rsmSdl), rsmdb.GetRsmWriter(rsmSdl))
	var msgImpl *rmrcgo.Context
	rmrMessenger := msgImpl.Init(config.Rmr.ReadyIntervalSec, "tcp:"+strconv.Itoa(config.Rmr.Port), config.Rmr.MaxMsgSize, 0, logger)
	rmrSender := rmrsender.NewRmrSender(logger, rmrMessenger)

	resourceStatusService := services.NewResourceStatusService(logger, rmrSender)
	unpacker := converters.NewX2apPduUnpacker(logger, e2pdus.MaxAsn1CodecMessageBufferSize)
	var rmrManager = rmrmanagers.NewRmrMessageManager(logger, config, rnibDataService, rmrSender, resourceStatusService, converters.NewResourceStatusResponseConverter(unpacker), converters.NewResourceStatusFailureConverter(unpacker))

	rmrReceiver := rmrreceiver.NewRmrReceiver(logger, rmrMessenger, rmrManager)
	defer rmrMessenger.Close()
	go rmrReceiver.ListenAndHandle()

	handlerProvider := httpmsghandlerprovider.NewRequestHandlerProvider(logger, rnibDataService, resourceStatusService)
	rootController := controllers.NewRootController(rnibDataService)
	controller := controllers.NewController(logger, handlerProvider)
	_ = httpserver.Run(config.Http.Port, rootController, controller)
}