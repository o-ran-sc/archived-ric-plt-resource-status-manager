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


package httpserver

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"rsm/mocks"
	"testing"
	"time"
)

func setupRouterAndMocks() (*mux.Router, *mocks.RootControllerMock, *mocks.ControllerMock) {
	rootControllerMock := &mocks.RootControllerMock{}
	rootControllerMock.On("HandleHealthCheckRequest").Return(nil)

	controllerMock := &mocks.ControllerMock{}
	controllerMock.On("ResourceStatus").Return(nil)

	router := mux.NewRouter()
	initializeRoutes(router, rootControllerMock, controllerMock)
	return router, rootControllerMock, controllerMock
}
func TestResourceStatus(t *testing.T) {
	router, _, controllerMock := setupRouterAndMocks()

	req, err := http.NewRequest("PUT", "/v1/general/resourcestatus", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	controllerMock.AssertNumberOfCalls(t, "ResourceStatus", 1)
}

func TestRouteGetHealth(t *testing.T) {
	router, rootControllerMock, _ := setupRouterAndMocks()

	req, err := http.NewRequest("GET", "/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	rootControllerMock.AssertNumberOfCalls(t, "HandleHealthCheckRequest", 1)
}

func TestRouteNotFound(t *testing.T) {
	router, _, _ := setupRouterAndMocks()

	req, err := http.NewRequest("GET", "/v1/no/such/route", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")
}

func TestRunError(t *testing.T) {
	_, rootControllerMock, controllerMock := setupRouterAndMocks()

	err := Run(111222333, rootControllerMock, controllerMock)

	assert.NotNil(t, err)
}

func TestRun(t *testing.T) {
	_, rootControllerMock, controllerMock := setupRouterAndMocks()

	go Run(11223, rootControllerMock, controllerMock)

	time.Sleep(time.Millisecond * 100)
	resp, err := http.Get("http://localhost:11223/v1/health")
	if err != nil {
		t.Fatalf("failed to perform GET to http://localhost:11223/v1/health")
	}
	assert.Equal(t, 200, resp.StatusCode)
}