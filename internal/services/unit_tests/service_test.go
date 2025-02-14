package unit_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hamillka/avitoTechWinter25/internal/repositories"
	"github.com/hamillka/avitoTechWinter25/internal/repositories/models"
	"github.com/hamillka/avitoTechWinter25/internal/services"
	"github.com/hamillka/avitoTechWinter25/internal/services/mocks"
	"github.com/ozontech/allure-go/pkg/framework/asserts_wrapper/require"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
)

type ServiceSuite struct {
	suite.Suite
}

type ServiceTestFixture struct {
	ctrl        *gomock.Controller
	userRepo    *mocks.MockUserRepository
	inventRepo  *mocks.MockInventoryRepository
	merchRepo   *mocks.MockMerchRepository
	txRepo      *mocks.MockTransactionRepository
	service     *services.AvitoShopService
	user        models.User
	userSecond  models.User
	inventory   models.Inventory
	merch       models.Merch
	transaction models.Transaction
}

func TestRunSuite(t *testing.T) {
	suite.RunSuite(t, new(ServiceSuite))
}

func NewDefaultUser(id int64, username, password string) models.User {
	return models.User{
		ID:       id,
		Username: username,
		Password: password,
		Coins:    1000,
	}
}

func NewDefaultInventory() models.Inventory {
	return models.Inventory{
		ID:      1,
		UserID:  1,
		MerchID: 1,
		Amount:  1,
	}
}

func NewDefaultMerch() models.Merch {
	return models.Merch{
		ID:   1,
		Type: "type",
		Cost: 10,
	}
}

func NewDefaultTransaction() models.Transaction {
	return models.Transaction{
		ID:         1,
		SenderID:   1,
		ReceiverID: 1,
		Amount:     50,
	}
}

func NewServiceTestFixture(t provider.T) *ServiceTestFixture {
	ctrl := gomock.NewController(t)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockInventRepo := mocks.NewMockInventoryRepository(ctrl)
	mockMerchRepo := mocks.NewMockMerchRepository(ctrl)
	mockTxRepo := mocks.NewMockTransactionRepository(ctrl)

	service := services.NewAvitoShopService(mockUserRepo, mockInventRepo, mockMerchRepo, mockTxRepo)

	return &ServiceTestFixture{
		ctrl:        ctrl,
		userRepo:    mockUserRepo,
		inventRepo:  mockInventRepo,
		merchRepo:   mockMerchRepo,
		txRepo:      mockTxRepo,
		service:     service,
		user:        NewDefaultUser(1, "test1", "test1"),
		userSecond:  NewDefaultUser(2, "test2", "test2"),
		inventory:   NewDefaultInventory(),
		merch:       NewDefaultMerch(),
		transaction: NewDefaultTransaction(),
	}
}

func (s *ServiceSuite) TestLoginExisting(t provider.T) {
	t.Title("User login if user exists in system")
	t.Tags("login", "success")

	fixture := NewServiceTestFixture(t)

	fixture.userRepo.EXPECT().
		GetUserByUsernamePassword(fixture.user.Username, fixture.user.Password).Return(fixture.user, nil)

	user, err := fixture.service.Login(fixture.user.Username, fixture.user.Password)

	require.NoError(t, err)
	assert.Equal(t, fixture.user.ID, user.ID)
}

func (s *ServiceSuite) TestLoginNotExisting(t provider.T) {
	t.Title("User login if user does not exist in system")
	t.Tags("login", "success")

	fixture := NewServiceTestFixture(t)

	fixture.userRepo.EXPECT().
		GetUserByUsernamePassword(fixture.user.Username, fixture.user.Password).
		Return(models.User{}, repositories.ErrRecordNotFound)
	fixture.userRepo.EXPECT().
		CreateUser(fixture.user.Username, fixture.user.Password).
		Return(fixture.user, nil)

	user, err := fixture.service.Login(fixture.user.Username, fixture.user.Password)

	require.NoError(t, err)
	assert.Equal(t, fixture.user.ID, user.ID)
}

func (s *ServiceSuite) TestBuyItemSuccess(t provider.T) {
	t.Title("Item buying success: merch and user exist")
	t.Tags("merch", "success")

	fixture := NewServiceTestFixture(t)

	fixture.merchRepo.EXPECT().
		GetMerchByType(fixture.merch.Type).
		Return(fixture.merch, nil)
	fixture.userRepo.EXPECT().
		GetUserByUsername(fixture.user.Username).
		Return(fixture.user, nil)
	fixture.userRepo.EXPECT().
		BuyItemFromAvitoShop(fixture.user.ID, fixture.merch.ID, fixture.merch.Cost).
		Return(nil)

	err := fixture.service.BuyItem(fixture.user.Username, fixture.merch.Type)

	require.NoError(t, err)
}

func (s *ServiceSuite) TestBuyItemMerchTypeError(t provider.T) {
	t.Title("Item buying error: merch does not exist")
	t.Tags("merch", "error")

	fixture := NewServiceTestFixture(t)

	fixture.merchRepo.EXPECT().
		GetMerchByType(fixture.merch.Type).
		Return(models.Merch{}, repositories.ErrRecordNotFound)

	err := fixture.service.BuyItem(fixture.user.Username, fixture.merch.Type)

	require.Error(t, err)
}

func (s *ServiceSuite) TestBuyItemUsernameError(t provider.T) {
	t.Title("Item buying error: user does not exist")
	t.Tags("merch", "error")

	fixture := NewServiceTestFixture(t)

	fixture.merchRepo.EXPECT().
		GetMerchByType(fixture.merch.Type).
		Return(fixture.merch, nil)
	fixture.userRepo.EXPECT().
		GetUserByUsername(fixture.user.Username).
		Return(models.User{}, repositories.ErrRecordNotFound)

	err := fixture.service.BuyItem(fixture.user.Username, fixture.merch.Type)

	require.Error(t, err)
}

func (s *ServiceSuite) TestBuyItemBuyingError(t provider.T) {
	t.Title("Item buying error: buying process error")
	t.Tags("merch", "error")

	fixture := NewServiceTestFixture(t)

	fixture.merchRepo.EXPECT().
		GetMerchByType(fixture.merch.Type).
		Return(fixture.merch, nil)
	fixture.userRepo.EXPECT().
		GetUserByUsername(fixture.user.Username).
		Return(fixture.user, nil)
	fixture.userRepo.EXPECT().
		BuyItemFromAvitoShop(fixture.user.ID, fixture.merch.ID, fixture.merch.Cost).
		Return(repositories.ErrDatabaseUpdatingError)

	err := fixture.service.BuyItem(fixture.user.Username, fixture.merch.Type)

	require.Error(t, err)
}

func (s *ServiceSuite) TestSendCoinSuccess(t provider.T) {
	t.Title("Coin sending success")
	t.Tags("coin", "success")

	fixture := NewServiceTestFixture(t)

	fixture.userRepo.EXPECT().
		GetUserByUsername(fixture.user.Username).
		Return(fixture.user, nil)
	fixture.userRepo.EXPECT().
		GetUserByUsername(fixture.userSecond.Username).
		Return(fixture.userSecond, nil)
	fixture.userRepo.EXPECT().
		TransferCoins(fixture.user.ID, fixture.userSecond.ID, fixture.transaction.Amount).
		Return(nil)

	err := fixture.service.SendCoin(fixture.user.Username, fixture.userSecond.Username, fixture.transaction.Amount)

	require.NoError(t, err)
}

func (s *ServiceSuite) TestSendCoinUsernameError(t provider.T) {
	t.Title("Coin sending error: first user does not exist")
	t.Tags("coin", "error")

	fixture := NewServiceTestFixture(t)

	fixture.userRepo.EXPECT().
		GetUserByUsername(fixture.user.Username).
		Return(models.User{}, repositories.ErrRecordNotFound)

	err := fixture.service.SendCoin(fixture.user.Username, fixture.userSecond.Username, fixture.transaction.Amount)

	require.Error(t, err)
}

func (s *ServiceSuite) TestSendCoinUsernameSecondError(t provider.T) {
	t.Title("Coin sending error: second user does not exist")
	t.Tags("coin", "error")

	fixture := NewServiceTestFixture(t)

	fixture.userRepo.EXPECT().
		GetUserByUsername(fixture.user.Username).
		Return(fixture.user, nil)
	fixture.userRepo.EXPECT().
		GetUserByUsername(fixture.userSecond.Username).
		Return(models.User{}, repositories.ErrRecordNotFound)

	err := fixture.service.SendCoin(fixture.user.Username, fixture.userSecond.Username, fixture.transaction.Amount)

	require.Error(t, err)
}

func (s *ServiceSuite) TestSendCoinTransferError(t provider.T) {
	t.Title("Coin sending error: transfer error")
	t.Tags("coin", "error")

	fixture := NewServiceTestFixture(t)

	fixture.userRepo.EXPECT().
		GetUserByUsername(fixture.user.Username).
		Return(fixture.user, nil)
	fixture.userRepo.EXPECT().
		GetUserByUsername(fixture.userSecond.Username).
		Return(fixture.userSecond, nil)
	fixture.userRepo.EXPECT().
		TransferCoins(fixture.user.ID, fixture.userSecond.ID, fixture.transaction.Amount).
		Return(repositories.ErrDatabaseUpdatingError)

	err := fixture.service.SendCoin(fixture.user.Username, fixture.userSecond.Username, fixture.transaction.Amount)

	require.Error(t, err)
}

func (s *ServiceSuite) TestGetInfoSuccess(t provider.T) {
	t.Title("Get info success")
	t.Tags("info", "success")

	fixture := NewServiceTestFixture(t)

	fixture.userRepo.EXPECT().
		GetUserByUsername(fixture.user.Username).
		Return(fixture.user, nil)
	fixture.inventRepo.EXPECT().
		GetInventoryByUserID(fixture.user.ID).
		Return([]*models.Inventory{&fixture.inventory}, nil)
	fixture.merchRepo.EXPECT().
		GetMerchByID(fixture.merch.ID).
		Return(fixture.merch, nil)
	fixture.txRepo.EXPECT().
		GetOutTransactions(fixture.user.ID).
		Return([]*models.Transaction{&fixture.transaction}, nil)
	fixture.userRepo.EXPECT().
		GetUserByID(fixture.user.ID).
		Return(fixture.user, nil)
	fixture.txRepo.EXPECT().
		GetInTransactions(fixture.user.ID).
		Return([]*models.Transaction{&fixture.transaction}, nil)
	fixture.userRepo.EXPECT().
		GetUserByID(fixture.user.ID).
		Return(fixture.user, nil)

	info, err := fixture.service.GetInfo(fixture.user.Username)

	require.NoError(t, err)
	assert.NotNil(t, info)
}

func (s *ServiceSuite) TestGetInfoUsernameError(t provider.T) {
	t.Title("Get info error: username not found")
	t.Tags("info", "error")

	fixture := NewServiceTestFixture(t)

	fixture.userRepo.EXPECT().
		GetUserByUsername(fixture.user.Username).
		Return(models.User{}, repositories.ErrRecordNotFound)

	info, err := fixture.service.GetInfo(fixture.user.Username)

	require.Error(t, err)
	assert.Nil(t, info)
}

func (s *ServiceSuite) TestGetInfoInventoryError(t provider.T) {
	t.Title("Get info error: inventory reading error")
	t.Tags("info", "error")

	fixture := NewServiceTestFixture(t)

	fixture.userRepo.EXPECT().
		GetUserByUsername(fixture.user.Username).
		Return(fixture.user, nil)
	fixture.inventRepo.EXPECT().
		GetInventoryByUserID(fixture.user.ID).
		Return([]*models.Inventory{}, repositories.ErrDatabaseReadingError)

	info, err := fixture.service.GetInfo(fixture.user.Username)

	require.Error(t, err)
	assert.Nil(t, info)
}

//func (s *ServiceSuite) TestGetInfoMerchError(t provider.T) {
//	t.Title("Get info error: merch not found")
//	t.Tags("info", "error")
//
//	fixture := NewServiceTestFixture(t)
//
//	fixture.userRepo.EXPECT().
//		GetUserByUsername(fixture.user.Username).
//		Return(fixture.user, nil)
//	fixture.inventRepo.EXPECT().
//		GetInventoryByUserID(fixture.user.ID).
//		Return([]*models.Inventory{&fixture.inventory}, nil)
//	fixture.merchRepo.EXPECT().
//		GetMerchByID(fixture.merch.ID).
//		Return(models.Merch{}, repositories.ErrRecordNotFound)
//
//	info, err := fixture.service.GetInfo(fixture.user.Username)
//
//	require.Error(t, err)
//	assert.Nil(t, info)
//}

func (s *ServiceSuite) TestGetInfoOutgoingTransactionError(t provider.T) {
	t.Title("Get info error: out transaction reading error")
	t.Tags("info", "error")

	fixture := NewServiceTestFixture(t)

	fixture.userRepo.EXPECT().
		GetUserByUsername(fixture.user.Username).
		Return(fixture.user, nil)
	fixture.inventRepo.EXPECT().
		GetInventoryByUserID(fixture.user.ID).
		Return([]*models.Inventory{&fixture.inventory}, nil)
	fixture.merchRepo.EXPECT().
		GetMerchByID(fixture.merch.ID).
		Return(fixture.merch, nil)
	fixture.txRepo.EXPECT().
		GetOutTransactions(fixture.user.ID).
		Return([]*models.Transaction{}, repositories.ErrDatabaseReadingError)

	info, err := fixture.service.GetInfo(fixture.user.Username)

	require.Error(t, err)
	assert.Nil(t, info)
}

//func (s *ServiceSuite) TestGetInfoGetUserError(t provider.T) {
//	t.Title("Get info error: get user by id error")
//	t.Tags("info", "error")
//
//	fixture := NewServiceTestFixture(t)
//
//	fixture.userRepo.EXPECT().
//		GetUserByUsername(fixture.user.Username).
//		Return(fixture.user, nil)
//	fixture.inventRepo.EXPECT().
//		GetInventoryByUserID(fixture.user.ID).
//		Return([]*models.Inventory{&fixture.inventory}, nil)
//	fixture.merchRepo.EXPECT().
//		GetMerchByID(fixture.merch.ID).
//		Return(fixture.merch, nil)
//	fixture.txRepo.EXPECT().
//		GetOutTransactions(fixture.user.ID).
//		Return([]*models.Transaction{&fixture.transaction}, nil)
//	fixture.userRepo.EXPECT().
//		GetUserByID(fixture.user.ID).
//		Return(models.User{}, repositories.ErrRecordNotFound)
//
//	info, err := fixture.service.GetInfo(fixture.user.Username)
//
//	require.Error(t, err)
//	assert.Nil(t, info)
//}

func (s *ServiceSuite) TestGetInfoIncomingTransactionError(t provider.T) {
	t.Title("Get info error: in transaction reading error")
	t.Tags("info", "error")

	fixture := NewServiceTestFixture(t)

	fixture.userRepo.EXPECT().
		GetUserByUsername(fixture.user.Username).
		Return(fixture.user, nil)
	fixture.inventRepo.EXPECT().
		GetInventoryByUserID(fixture.user.ID).
		Return([]*models.Inventory{&fixture.inventory}, nil)
	fixture.merchRepo.EXPECT().
		GetMerchByID(fixture.merch.ID).
		Return(fixture.merch, nil)
	fixture.txRepo.EXPECT().
		GetOutTransactions(fixture.user.ID).
		Return([]*models.Transaction{&fixture.transaction}, nil)
	fixture.userRepo.EXPECT().
		GetUserByID(fixture.user.ID).
		Return(fixture.user, nil)
	fixture.txRepo.EXPECT().
		GetInTransactions(fixture.user.ID).
		Return([]*models.Transaction{}, repositories.ErrDatabaseReadingError)

	info, err := fixture.service.GetInfo(fixture.user.Username)

	require.Error(t, err)
	assert.Nil(t, info)
}
