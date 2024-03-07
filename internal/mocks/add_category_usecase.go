// Code generated by MockGen. DO NOT EDIT.
// Source: internal/controller/http/v1/handler/category/add.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/The-Gleb/product_catalog/internal/domain/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockAddCategoryUsecase is a mock of AddCategoryUsecase interface.
type MockAddCategoryUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockAddCategoryUsecaseMockRecorder
}

// MockAddCategoryUsecaseMockRecorder is the mock recorder for MockAddCategoryUsecase.
type MockAddCategoryUsecaseMockRecorder struct {
	mock *MockAddCategoryUsecase
}

// NewMockAddCategoryUsecase creates a new mock instance.
func NewMockAddCategoryUsecase(ctrl *gomock.Controller) *MockAddCategoryUsecase {
	mock := &MockAddCategoryUsecase{ctrl: ctrl}
	mock.recorder = &MockAddCategoryUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAddCategoryUsecase) EXPECT() *MockAddCategoryUsecaseMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockAddCategoryUsecase) Add(ctx context.Context, category entity.AddCategoryDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, category)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockAddCategoryUsecaseMockRecorder) Add(ctx, category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockAddCategoryUsecase)(nil).Add), ctx, category)
}