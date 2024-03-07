// Code generated by MockGen. DO NOT EDIT.
// Source: internal/controller/http/v1/handler/product/update_category.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entity "github.com/The-Gleb/product_catalog/internal/domain/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockUpdateProductCategoryUsecase is a mock of UpdateProductCategoryUsecase interface.
type MockUpdateProductCategoryUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUpdateProductCategoryUsecaseMockRecorder
}

// MockUpdateProductCategoryUsecaseMockRecorder is the mock recorder for MockUpdateProductCategoryUsecase.
type MockUpdateProductCategoryUsecaseMockRecorder struct {
	mock *MockUpdateProductCategoryUsecase
}

// NewMockUpdateProductCategoryUsecase creates a new mock instance.
func NewMockUpdateProductCategoryUsecase(ctrl *gomock.Controller) *MockUpdateProductCategoryUsecase {
	mock := &MockUpdateProductCategoryUsecase{ctrl: ctrl}
	mock.recorder = &MockUpdateProductCategoryUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUpdateProductCategoryUsecase) EXPECT() *MockUpdateProductCategoryUsecaseMockRecorder {
	return m.recorder
}

// UpdateCategory mocks base method.
func (m *MockUpdateProductCategoryUsecase) UpdateCategory(ctx context.Context, product entity.UpdateProductCategoryDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCategory", ctx, product)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCategory indicates an expected call of UpdateCategory.
func (mr *MockUpdateProductCategoryUsecaseMockRecorder) UpdateCategory(ctx, product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCategory", reflect.TypeOf((*MockUpdateProductCategoryUsecase)(nil).UpdateCategory), ctx, product)
}
