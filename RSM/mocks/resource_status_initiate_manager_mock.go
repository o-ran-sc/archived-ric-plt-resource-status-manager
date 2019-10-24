package mocks

import (
	"github.com/stretchr/testify/mock"
	"rsm/e2pdus"
)

type ResourceStatusInitiateManagerMock struct {
	mock.Mock
}

func (m *ResourceStatusInitiateManagerMock) Execute(inventoryName string, resourceStatusInitiateRequestParams *e2pdus.ResourceStatusRequestData) error {
	args := m.Called(inventoryName, resourceStatusInitiateRequestParams)
	return args.Error(0)
}
