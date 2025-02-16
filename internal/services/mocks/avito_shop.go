// Code generated by MockGen. DO NOT EDIT.
// Source: avito_shop.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/hamillka/avitoTechWinter25/internal/repositories/models"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// BuyItemFromAvitoShop mocks base method.
func (m *MockUserRepository) BuyItemFromAvitoShop(buyerID, itemID, itemCost int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuyItemFromAvitoShop", buyerID, itemID, itemCost)
	ret0, _ := ret[0].(error)
	return ret0
}

// BuyItemFromAvitoShop indicates an expected call of BuyItemFromAvitoShop.
func (mr *MockUserRepositoryMockRecorder) BuyItemFromAvitoShop(buyerID, itemID, itemCost interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuyItemFromAvitoShop", reflect.TypeOf((*MockUserRepository)(nil).BuyItemFromAvitoShop), buyerID, itemID, itemCost)
}

// CreateUser mocks base method.
func (m *MockUserRepository) CreateUser(username, password string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", username, password)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepositoryMockRecorder) CreateUser(username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), username, password)
}

// GetUserByID mocks base method.
func (m *MockUserRepository) GetUserByID(id int64) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", id)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserRepositoryMockRecorder) GetUserByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserRepository)(nil).GetUserByID), id)
}

// GetUserByUsername mocks base method.
func (m *MockUserRepository) GetUserByUsername(username string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", username)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockUserRepositoryMockRecorder) GetUserByUsername(username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockUserRepository)(nil).GetUserByUsername), username)
}

// GetUserByUsernamePassword mocks base method.
func (m *MockUserRepository) GetUserByUsernamePassword(username, password string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsernamePassword", username, password)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsernamePassword indicates an expected call of GetUserByUsernamePassword.
func (mr *MockUserRepositoryMockRecorder) GetUserByUsernamePassword(username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsernamePassword", reflect.TypeOf((*MockUserRepository)(nil).GetUserByUsernamePassword), username, password)
}

// TransferCoins mocks base method.
func (m *MockUserRepository) TransferCoins(senderID, receiverID, amount int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferCoins", senderID, receiverID, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// TransferCoins indicates an expected call of TransferCoins.
func (mr *MockUserRepositoryMockRecorder) TransferCoins(senderID, receiverID, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferCoins", reflect.TypeOf((*MockUserRepository)(nil).TransferCoins), senderID, receiverID, amount)
}

// MockInventoryRepository is a mock of InventoryRepository interface.
type MockInventoryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockInventoryRepositoryMockRecorder
}

// MockInventoryRepositoryMockRecorder is the mock recorder for MockInventoryRepository.
type MockInventoryRepositoryMockRecorder struct {
	mock *MockInventoryRepository
}

// NewMockInventoryRepository creates a new mock instance.
func NewMockInventoryRepository(ctrl *gomock.Controller) *MockInventoryRepository {
	mock := &MockInventoryRepository{ctrl: ctrl}
	mock.recorder = &MockInventoryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInventoryRepository) EXPECT() *MockInventoryRepositoryMockRecorder {
	return m.recorder
}

// GetInventoryByUserID mocks base method.
func (m *MockInventoryRepository) GetInventoryByUserID(userID int64) ([]*models.Inventory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInventoryByUserID", userID)
	ret0, _ := ret[0].([]*models.Inventory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInventoryByUserID indicates an expected call of GetInventoryByUserID.
func (mr *MockInventoryRepositoryMockRecorder) GetInventoryByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInventoryByUserID", reflect.TypeOf((*MockInventoryRepository)(nil).GetInventoryByUserID), userID)
}

// MockMerchRepository is a mock of MerchRepository interface.
type MockMerchRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMerchRepositoryMockRecorder
}

// MockMerchRepositoryMockRecorder is the mock recorder for MockMerchRepository.
type MockMerchRepositoryMockRecorder struct {
	mock *MockMerchRepository
}

// NewMockMerchRepository creates a new mock instance.
func NewMockMerchRepository(ctrl *gomock.Controller) *MockMerchRepository {
	mock := &MockMerchRepository{ctrl: ctrl}
	mock.recorder = &MockMerchRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMerchRepository) EXPECT() *MockMerchRepositoryMockRecorder {
	return m.recorder
}

// GetMerchByID mocks base method.
func (m *MockMerchRepository) GetMerchByID(id int64) (models.Merch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMerchByID", id)
	ret0, _ := ret[0].(models.Merch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMerchByID indicates an expected call of GetMerchByID.
func (mr *MockMerchRepositoryMockRecorder) GetMerchByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMerchByID", reflect.TypeOf((*MockMerchRepository)(nil).GetMerchByID), id)
}

// GetMerchByType mocks base method.
func (m *MockMerchRepository) GetMerchByType(merchType string) (models.Merch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMerchByType", merchType)
	ret0, _ := ret[0].(models.Merch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMerchByType indicates an expected call of GetMerchByType.
func (mr *MockMerchRepositoryMockRecorder) GetMerchByType(merchType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMerchByType", reflect.TypeOf((*MockMerchRepository)(nil).GetMerchByType), merchType)
}

// MockTransactionRepository is a mock of TransactionRepository interface.
type MockTransactionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionRepositoryMockRecorder
}

// MockTransactionRepositoryMockRecorder is the mock recorder for MockTransactionRepository.
type MockTransactionRepositoryMockRecorder struct {
	mock *MockTransactionRepository
}

// NewMockTransactionRepository creates a new mock instance.
func NewMockTransactionRepository(ctrl *gomock.Controller) *MockTransactionRepository {
	mock := &MockTransactionRepository{ctrl: ctrl}
	mock.recorder = &MockTransactionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionRepository) EXPECT() *MockTransactionRepositoryMockRecorder {
	return m.recorder
}

// GetInTransactions mocks base method.
func (m *MockTransactionRepository) GetInTransactions(userID int64) ([]*models.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInTransactions", userID)
	ret0, _ := ret[0].([]*models.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInTransactions indicates an expected call of GetInTransactions.
func (mr *MockTransactionRepositoryMockRecorder) GetInTransactions(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInTransactions", reflect.TypeOf((*MockTransactionRepository)(nil).GetInTransactions), userID)
}

// GetOutTransactions mocks base method.
func (m *MockTransactionRepository) GetOutTransactions(userID int64) ([]*models.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOutTransactions", userID)
	ret0, _ := ret[0].([]*models.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOutTransactions indicates an expected call of GetOutTransactions.
func (mr *MockTransactionRepositoryMockRecorder) GetOutTransactions(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOutTransactions", reflect.TypeOf((*MockTransactionRepository)(nil).GetOutTransactions), userID)
}
