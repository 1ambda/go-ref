// Code generated by MockGen. DO NOT EDIT.
// Source: internal/distributed/client.go

// Package distributed is a generated GoMock package.
package distributed

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockConnector is a mock of Connector interface
type MockConnector struct {
	ctrl     *gomock.Controller
	recorder *MockConnectorMockRecorder
}

// MockConnectorMockRecorder is the mock recorder for MockConnector
type MockConnectorMockRecorder struct {
	mock *MockConnector
}

// NewMockConnector creates a new mock instance
func NewMockConnector(ctrl *gomock.Controller) *MockConnector {
	mock := &MockConnector{ctrl: ctrl}
	mock.recorder = &MockConnectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConnector) EXPECT() *MockConnectorMockRecorder {
	return m.recorder
}

// GetLeaderOrCampaign mocks base method
func (m *MockConnector) GetLeaderOrCampaign(electSubPath, electProclaim string) (string, error) {
	ret := m.ctrl.Call(m, "GetLeaderOrCampaign", electSubPath, electProclaim)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLeaderOrCampaign indicates an expected call of GetLeaderOrCampaign
func (mr *MockConnectorMockRecorder) GetLeaderOrCampaign(electSubPath, electProclaim interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLeaderOrCampaign", reflect.TypeOf((*MockConnector)(nil).GetLeaderOrCampaign), electSubPath, electProclaim)
}

// Publish mocks base method
func (m *MockConnector) Publish(ctx context.Context, message *Message) error {
	ret := m.ctrl.Call(m, "Publish", ctx, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish
func (mr *MockConnectorMockRecorder) Publish(ctx, message interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockConnector)(nil).Publish), ctx, message)
}

// Stop mocks base method
func (m *MockConnector) Stop() {
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop
func (mr *MockConnectorMockRecorder) Stop() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockConnector)(nil).Stop))
}
