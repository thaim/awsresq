// Code generated by MockGen. DO NOT EDIT.
// Source: s3.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	gomock "github.com/golang/mock/gomock"
)

// MockawsS3API is a mock of awsS3API interface.
type MockawsS3API struct {
	ctrl     *gomock.Controller
	recorder *MockawsS3APIMockRecorder
}

// MockawsS3APIMockRecorder is the mock recorder for MockawsS3API.
type MockawsS3APIMockRecorder struct {
	mock *MockawsS3API
}

// NewMockawsS3API creates a new mock instance.
func NewMockawsS3API(ctrl *gomock.Controller) *MockawsS3API {
	mock := &MockawsS3API{ctrl: ctrl}
	mock.recorder = &MockawsS3APIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockawsS3API) EXPECT() *MockawsS3APIMockRecorder {
	return m.recorder
}

// ListBuckets mocks base method.
func (m *MockawsS3API) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListBuckets", varargs...)
	ret0, _ := ret[0].(*s3.ListBucketsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListBuckets indicates an expected call of ListBuckets.
func (mr *MockawsS3APIMockRecorder) ListBuckets(ctx, params interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBuckets", reflect.TypeOf((*MockawsS3API)(nil).ListBuckets), varargs...)
}
