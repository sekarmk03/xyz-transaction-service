package service_test

import (
	"context"
	// "strings"
	"testing"
	"time"
	"xyz-transaction-service/common/config"
	"xyz-transaction-service/modules/transaction/entity"
	"xyz-transaction-service/modules/transaction/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for TransactionRepositoryUseCase
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) FindAll(ctx context.Context, req any) ([]*entity.Transaction, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByConsumerId(ctx context.Context, consumerId uint64) ([]*entity.Transaction, error) {
	args := m.Called(ctx, consumerId)
	return args.Get(0).([]*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindById(ctx context.Context, id uint64) (*entity.Transaction, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByContractNumber(ctx context.Context, contractNumber string) (*entity.Transaction, error) {
	args := m.Called(ctx, contractNumber)
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error) {
	args := m.Called(ctx, transaction)
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestFindById(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	mockTransaction := &entity.Transaction{
		Id:             1,
		ConsumerId:     123,
		ContractNumber: "CN123",
		Tenor:          12,
		Otr:            300000,
		AdminFee:       20000,
		Installment:    10000,
		Interest:       5000,
		AssetName:      "Laptop",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mockRepo.On("FindById", mock.Anything, uint64(1)).Return(mockTransaction, nil)

	svc := service.NewTransactionService(config.Config{}, mockRepo)

	result, err := svc.FindById(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockTransaction.Id, result.Id)
	assert.Equal(t, mockTransaction.ConsumerId, result.ConsumerId)
	assert.Equal(t, "CN123", result.ContractNumber)

	mockRepo.AssertExpectations(t)
}

// func TestCreate(t *testing.T) {
// 	mockRepo := new(MockTransactionRepository)

// 	mockTransaction := &entity.Transaction{
// 		Id:             1,
// 		ContractNumber: "CN123",
// 		ConsumerId:     3,
// 		Tenor:          12,
// 		Otr:            300000,
// 		AdminFee:       18000,
// 		Installment:    135000,
// 		Interest:       12000,
// 		AssetName:      "Smartwatch",
// 		CreatedAt:      time.Now(),
// 		UpdatedAt:      time.Now(),
// 	}

// 	mockRepo.On("Create", mock.Anything, mockTransaction).Return(mockTransaction, nil)

// 	svc := service.NewTransactionService(config.Config{}, mockRepo)

// 	result, err := svc.Create(context.Background(), 3, 12, 300000, 18000, 135000, 12000, "Smartwatch")

// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(t, uint64(3), result.ConsumerId)
// 	assert.True(t, strings.HasPrefix(result.ContractNumber, "CNTR-"))
// 	assert.Equal(t, uint32(12), result.Tenor)
// 	assert.Equal(t, uint64(300000), result.Otr)
// 	assert.Equal(t, uint64(18000), result.AdminFee)
// 	assert.Equal(t, uint64(135000), result.Installment)
// 	assert.Equal(t, uint64(12000), result.Interest)
// 	assert.Equal(t, "Smartwatch", result.AssetName)
// 	assert.WithinDuration(t, mockTransaction.CreatedAt, result.CreatedAt, time.Second)
// 	assert.WithinDuration(t, mockTransaction.UpdatedAt, result.UpdatedAt, time.Second)

// 	mockRepo.AssertExpectations(t)
// }

func TestRollback(t *testing.T) {
	mockRepo := new(MockTransactionRepository)

	mockRepo.On("Delete", mock.Anything, uint64(1)).Return(nil)

	svc := service.NewTransactionService(config.Config{}, mockRepo)

	err := svc.Rollback(context.Background(), 1)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
