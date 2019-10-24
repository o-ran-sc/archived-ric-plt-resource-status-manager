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

package controllers

type IController interface {
}
/*
import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"rsm/logger"
	"rsm/models"
	"rsm/providers/httpmsghandlerprovider"
	"rsm/rsmerrors"
	"strings"
)

const (
	LimitRequest = 2000
)

type IController interface {
}

type Controller struct {
	logger          *logger.Logger
	handlerProvider *httpmsghandlerprovider.RequestHandlerProvider
}

func NewController(logger *logger.Logger, handlerProvider *httpmsghandlerprovider.RequestHandlerProvider) *Controller {
	return &Controller{
		logger:          logger,
		handlerProvider: handlerProvider,
	}
}

func (c *Controller) extractJsonBody(r *http.Request, request models.Request, writer http.ResponseWriter) bool {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, LimitRequest))

	if err != nil {
		c.logger.Errorf("[Client -> RSM] #Controller.extractJsonBody - unable to extract json body - error: %s", err)
		c.handleErrorResponse(rsmerrors.NewInvalidJsonError(), writer)
		return false
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		c.logger.Errorf("[Client -> RSM] #Controller.extractJsonBody - unable to extract json body - error: %s", err)
		c.handleErrorResponse(rsmerrors.NewInvalidJsonError(), writer)
		return false
	}

	return true
}

func (c *Controller) handleRequest(writer http.ResponseWriter, header *http.Header, requestName httpmsghandlerprovider.IncomingRequest, request models.Request, validateHeader bool) {

	if validateHeader {

		err := c.validateRequestHeader(header)
		if err != nil {
			c.handleErrorResponse(err, writer)
			return
		}
	}

	handler, err := c.handlerProvider.GetHandler(requestName)

	if err != nil {
		c.handleErrorResponse(err, writer)
		return
	}

	response, err := (*handler).Handle(request)

	if err != nil {
		c.handleErrorResponse(err, writer)
		return
	}

	if response == nil {
		writer.WriteHeader(http.StatusNoContent)
		c.logger.Infof("[RSM -> Client] #Controller.handleRequest - status response: %v", http.StatusNoContent)
		return
	}

	result, err := response.Marshal()

	if err != nil {
		c.handleErrorResponse(err, writer)
		return
	}

	c.logger.Infof("[RSM -> Client] #Controller.handleRequest - response: %s", result)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte(result))
}

func (c *Controller) validateRequestHeader(header *http.Header) error {

	if header.Get("Content-Type") != "application/json" {
		c.logger.Errorf("#Controller.validateRequestHeader - validation failure, incorrect content type")

		return rsmerrors.NewHeaderValidationError()
	}
	return nil
}

func (c *Controller) handleErrorResponse(err error, writer http.ResponseWriter) {

	var errorResponseDetails models.ErrorResponse
	var httpError int

	if err != nil {
		switch err.(type) {
		case *rsmerrors.RnibDbError:
			e2Error, _ := err.(*rsmerrors.RnibDbError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusInternalServerError
		case *rsmerrors.HeaderValidationError:
			e2Error, _ := err.(*rsmerrors.HeaderValidationError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusUnsupportedMediaType
		case *rsmerrors.WrongStateError:
			e2Error, _ := err.(*rsmerrors.WrongStateError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusBadRequest
		case *rsmerrors.RequestValidationError:
			e2Error, _ := err.(*rsmerrors.RequestValidationError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusBadRequest
		case *rsmerrors.InvalidJsonError:
			e2Error, _ := err.(*rsmerrors.InvalidJsonError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusBadRequest
		case *rsmerrors.RmrError:
			e2Error, _ := err.(*rsmerrors.RmrError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusInternalServerError
		case *rsmerrors.ResourceNotFoundError:
			e2Error, _ := err.(*rsmerrors.ResourceNotFoundError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusNotFound

		default:
			e2Error := rsmerrors.NewInternalError()
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusInternalServerError
		}
	}
	errorResponse, _ := json.Marshal(errorResponseDetails)

	c.logger.Errorf("[RSM -> Client] #Controller.handleErrorResponse - http status: %d, error response: %+v", httpError, errorResponseDetails)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpError)
	_, err = writer.Write(errorResponse)

	if err != nil {
		c.logger.Errorf("#Controller.handleErrorResponse - Cannot send response. writer:%v", writer)
	}
}*/
/*
func (c *Controller) prettifyRequest(request *http.Request) string {
	dump, _ := httputil.DumpRequest(request, true)
	requestPrettyPrint := strings.Replace(string(dump), "\r\n", " ", -1)
	return strings.Replace(requestPrettyPrint, "\n", "", -1)
}
*/