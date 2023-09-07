// Code generated by MockGen. DO NOT EDIT.
// Source: efs.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	efs "github.com/aws/aws-sdk-go-v2/service/efs"
	gomock "github.com/golang/mock/gomock"
)

// MockawsEfsAPI is a mock of awsEfsAPI interface.
type MockawsEfsAPI struct {
	ctrl     *gomock.Controller
	recorder *MockawsEfsAPIMockRecorder
}

// MockawsEfsAPIMockRecorder is the mock recorder for MockawsEfsAPI.
type MockawsEfsAPIMockRecorder struct {
	mock *MockawsEfsAPI
}

// NewMockawsEfsAPI creates a new mock instance.
func NewMockawsEfsAPI(ctrl *gomock.Controller) *MockawsEfsAPI {
	mock := &MockawsEfsAPI{ctrl: ctrl}
	mock.recorder = &MockawsEfsAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockawsEfsAPI) EXPECT() *MockawsEfsAPIMockRecorder {
	return m.recorder
}

// DescribeFileSystems mocks base method.
func (m *MockawsEfsAPI) DescribeFileSystems(ctx context.Context, params *efs.DescribeFileSystemsInput, optFns ...func(*efs.Options)) (*efs.DescribeFileSystemsOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeFileSystems", varargs...)
	ret0, _ := ret[0].(*efs.DescribeFileSystemsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeFileSystems indicates an expected call of DescribeFileSystems.
func (mr *MockawsEfsAPIMockRecorder) DescribeFileSystems(ctx, params interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeFileSystems", reflect.TypeOf((*MockawsEfsAPI)(nil).DescribeFileSystems), varargs...)
}
