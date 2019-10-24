package mocks

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type RootControllerMock struct {
	mock.Mock
}

func (rc *RootControllerMock) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	rc.Called()
}

func (rc *RootControllerMock) HandleHealthCheckRequest(writer http.ResponseWriter, request *http.Request) {
	rc.Called()
}
