package service

import (
	"context"
	"log"
	"time"
	"xyz-transaction-service/common/config"
	commonErr "xyz-transaction-service/common/error"
	"xyz-transaction-service/common/utils"
	"xyz-transaction-service/modules/transaction/entity"
	"xyz-transaction-service/modules/transaction/internal/repository"
)

type TransactionService struct {
	cfg                   config.Config
	transactionRepository repository.TransactionRepositoryUseCase
}

func NewTransactionService(cfg config.Config, transactionRepository repository.TransactionRepositoryUseCase) *TransactionService {
	return &TransactionService{
		cfg:                   cfg,
		transactionRepository: transactionRepository,
	}
}

type TransactionServiceUseCase interface {
	FindAll(ctx context.Context, req any) ([]*entity.Transaction, error)
	FindByConsumerId(ctx context.Context, consumerId uint64) ([]*entity.Transaction, error)
	FindById(ctx context.Context, id uint64) (*entity.Transaction, error)
	FindByContractNumber(ctx context.Context, contractNumber string) (*entity.Transaction, error)
	Create(ctx context.Context, consumerId uint64, tenor uint32, otr, adminFee, installment, interest uint64, assetName string) (*entity.Transaction, error)
	Rollback(ctx context.Context, id uint64) error
}

func (svc *TransactionService) FindAll(ctx context.Context, req any) ([]*entity.Transaction, error) {
	res, err := svc.transactionRepository.FindAll(ctx, req)
	if err != nil {
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [TransactionService - FindAll] Error while find all transaction:", parseError.Message)
		return nil, err
	}

	return res, nil
}

func (svc *TransactionService) FindByConsumerId(ctx context.Context, consumerId uint64) ([]*entity.Transaction, error) {
	res, err := svc.transactionRepository.FindByConsumerId(ctx, consumerId)
	if err != nil {
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [TransactionService - FindByConsumerId] Error while find transaction by consumer id:", parseError.Message)
		return nil, err
	}

	return res, nil
}

func (svc *TransactionService) FindById(ctx context.Context, id uint64) (*entity.Transaction, error) {
	res, err := svc.transactionRepository.FindById(ctx, id)
	if err != nil {
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [TransactionService - FindById] Error while find transaction by id:", parseError.Message)
		return nil, err
	}

	return res, nil
}

func (svc *TransactionService) FindByContractNumber(ctx context.Context, contractNumber string) (*entity.Transaction, error) {
	res, err := svc.transactionRepository.FindByContractNumber(ctx, contractNumber)
	if err != nil {
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [TransactionService - FindByContractNumber] Error while find transaction by contract number:", parseError.Message)
		return nil, err
	}

	return res, nil
}

func (svc *TransactionService) Create(ctx context.Context, consumerId uint64, tenor uint32, otr, adminFee, installment, interest uint64, assetName string) (*entity.Transaction, error) {
	transaction := &entity.Transaction{
		ContractNumber: utils.GenerateContractNumber(consumerId),
		ConsumerId:     consumerId,
		Tenor:          tenor,
		Otr:            otr,
		AdminFee:       adminFee,
		Installment:    installment,
		Interest:       interest,
		AssetName:      assetName,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	res, err := svc.transactionRepository.Create(ctx, transaction)
	if err != nil {
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [TransactionService - Create] Error while create transaction:", parseError.Message)
		return nil, err
	}

	return res, nil
}

func (svc *TransactionService) Rollback(ctx context.Context, id uint64) error {
	err := svc.transactionRepository.Delete(ctx, id)
	if err != nil {
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [TransactionService - Rollback] Error while rollback transaction:", parseError.Message)
		return err
	}

	return nil
}
