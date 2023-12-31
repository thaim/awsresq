// Code generated by MockGen. DO NOT EDIT.
// Source: logs.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	cloudwatchlogs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	gomock "github.com/golang/mock/gomock"
)

// MockawsLogsAPI is a mock of awsLogsAPI interface.
type MockawsLogsAPI struct {
	ctrl     *gomock.Controller
	recorder *MockawsLogsAPIMockRecorder
}

// MockawsLogsAPIMockRecorder is the mock recorder for MockawsLogsAPI.
type MockawsLogsAPIMockRecorder struct {
	mock *MockawsLogsAPI
}

// NewMockawsLogsAPI creates a new mock instance.
func NewMockawsLogsAPI(ctrl *gomock.Controller) *MockawsLogsAPI {
	mock := &MockawsLogsAPI{ctrl: ctrl}
	mock.recorder = &MockawsLogsAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockawsLogsAPI) EXPECT() *MockawsLogsAPIMockRecorder {
	return m.recorder
}

// DescribeLogGroups mocks base method.
func (m *MockawsLogsAPI) DescribeLogGroups(ctx context.Context, params *cloudwatchlogs.DescribeLogGroupsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeLogGroups", varargs...)
	ret0, _ := ret[0].(*cloudwatchlogs.DescribeLogGroupsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeLogGroups indicates an expected call of DescribeLogGroups.
func (mr *MockawsLogsAPIMockRecorder) DescribeLogGroups(ctx, params interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeLogGroups", reflect.TypeOf((*MockawsLogsAPI)(nil).DescribeLogGroups), varargs...)
}
