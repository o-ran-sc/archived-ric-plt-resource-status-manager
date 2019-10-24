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

package main

import (
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"os"
	"rsm/configuration"
	"rsm/controllers"
	"rsm/converters"
	"rsm/httpserver"
	"rsm/logger"
	"rsm/managers"
	"rsm/managers/rmrmanagers"
	"rsm/rmrcgo"
	"rsm/services"
	"rsm/services/rmrreceiver"
	"rsm/services/rmrsender"
	"strconv"
)

const MaxRnibPoolInstances = 4

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
	reader.Init("e2Manager", MaxRnibPoolInstances)
	logger.Infof("#app.main - Configuration %s", config)
	defer reader.Close()
	rnibDataService := services.NewRnibDataService(logger, config, reader.GetRNibReader)
	var msgImpl *rmrcgo.Context
	rmrMessenger := msgImpl.Init(config.Rmr.ReadyIntervalSec, "tcp:"+strconv.Itoa(config.Rmr.Port), config.Rmr.MaxMsgSize, 0, logger)
	rmrSender := rmrsender.NewRmrSender(logger, rmrMessenger)

	resourceStatusInitiateManager := managers.NewResourceStatusInitiateManager(logger, rnibDataService, rmrSender)
	var rmrManager = rmrmanagers.NewRmrMessageManager(logger, config, rnibDataService, rmrSender, resourceStatusInitiateManager,converters.NewX2apPduUnpacker(logger))

	rmrReceiver := rmrreceiver.NewRmrReceiver(logger, rmrMessenger, rmrManager)
	defer rmrMessenger.Close()
	go rmrReceiver.ListenAndHandle()

	//handlerProvider := httpmsghandlerprovider.NewRequestHandlerProvider(logger, rmrSender, config, rnibDataService)
	rootController := controllers.NewRootController(rnibDataService)
	//controller := controllers.NewController(logger, handlerProvider)
	_ = httpserver.Run(config.Http.Port, rootController)
}
