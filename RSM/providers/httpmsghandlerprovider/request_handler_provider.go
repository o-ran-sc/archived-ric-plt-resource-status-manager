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

package httpmsghandlerprovider

import (
	"rsm/handlers/httpmsghandlers"
	"rsm/logger"
	"rsm/rsmerrors"
	"rsm/services"
)

type IncomingRequest string

const (
	ResourceStatusRequest = "ResourceStatusRequest"
)

type RequestHandlerProvider struct {
	requestMap map[IncomingRequest]httpmsghandlers.RequestHandler
	logger     *logger.Logger
}

func NewRequestHandlerProvider(logger *logger.Logger, rNibDataService services.RNibDataService, resourceStatusService services.IResourceStatusService) *RequestHandlerProvider {

	return &RequestHandlerProvider{
		requestMap: initRequestHandlerMap(logger, rNibDataService, resourceStatusService),
		logger:     logger,
	}
}

func initRequestHandlerMap(logger *logger.Logger, rNibDataService services.RNibDataService, resourceStatusService services.IResourceStatusService) map[IncomingRequest]httpmsghandlers.RequestHandler {

	return map[IncomingRequest]httpmsghandlers.RequestHandler{
		ResourceStatusRequest: httpmsghandlers.NewResourceStatusRequestHandler(logger, rNibDataService, resourceStatusService),
	}
}

func (provider RequestHandlerProvider) GetHandler(requestType IncomingRequest) (httpmsghandlers.RequestHandler, error) {
	handler, ok := provider.requestMap[requestType]

	if !ok {
		provider.logger.Errorf("#request_handler_provider.GetHandler - Cannot find handler for request type: %s", requestType)
		return nil, rsmerrors.NewInternalError()
	}

	return handler, nil
}