// Code generated by MockGen. DO NOT EDIT.
// Source: ec2.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	ec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	gomock "github.com/golang/mock/gomock"
)

// MockawsEc2API is a mock of awsEc2API interface.
type MockawsEc2API struct {
	ctrl     *gomock.Controller
	recorder *MockawsEc2APIMockRecorder
}

// MockawsEc2APIMockRecorder is the mock recorder for MockawsEc2API.
type MockawsEc2APIMockRecorder struct {
	mock *MockawsEc2API
}

// NewMockawsEc2API creates a new mock instance.
func NewMockawsEc2API(ctrl *gomock.Controller) *MockawsEc2API {
	mock := &MockawsEc2API{ctrl: ctrl}
	mock.recorder = &MockawsEc2APIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockawsEc2API) EXPECT() *MockawsEc2APIMockRecorder {
	return m.recorder
}

// DescribeInstances mocks base method.
func (m *MockawsEc2API) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeInstances", varargs...)
	ret0, _ := ret[0].(*ec2.DescribeInstancesOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeInstances indicates an expected call of DescribeInstances.
func (mr *MockawsEc2APIMockRecorder) DescribeInstances(ctx, params interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeInstances", reflect.TypeOf((*MockawsEc2API)(nil).DescribeInstances), varargs...)
}

// DescribeSecurityGroups mocks base method.
func (m *MockawsEc2API) DescribeSecurityGroups(ctx context.Context, params *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeSecurityGroups", varargs...)
	ret0, _ := ret[0].(*ec2.DescribeSecurityGroupsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeSecurityGroups indicates an expected call of DescribeSecurityGroups.
func (mr *MockawsEc2APIMockRecorder) DescribeSecurityGroups(ctx, params interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeSecurityGroups", reflect.TypeOf((*MockawsEc2API)(nil).DescribeSecurityGroups), varargs...)
}

// DescribeVpcs mocks base method.
func (m *MockawsEc2API) DescribeVpcs(ctx context.Context, params *ec2.DescribeVpcsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeVpcs", varargs...)
	ret0, _ := ret[0].(*ec2.DescribeVpcsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeVpcs indicates an expected call of DescribeVpcs.
func (mr *MockawsEc2APIMockRecorder) DescribeVpcs(ctx, params interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeVpcs", reflect.TypeOf((*MockawsEc2API)(nil).DescribeVpcs), varargs...)
}
