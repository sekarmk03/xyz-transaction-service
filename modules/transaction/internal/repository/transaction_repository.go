package repository

import (
	"context"
	"errors"
	"log"
	"xyz-transaction-service/modules/transaction/entity"

	"github.com/go-sql-driver/mysql"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

type TransactionRepositoryUseCase interface {
	FindAll(ctx context.Context, req any) ([]*entity.Transaction, error)
	FindByConsumerId(ctx context.Context, consumerId uint64) ([]*entity.Transaction, error)
	FindById(ctx context.Context, id uint64) (*entity.Transaction, error)
	FindByContractNumber(ctx context.Context, contractNumber string) (*entity.Transaction, error)
	Create(ctx context.Context, req *entity.Transaction) (*entity.Transaction, error)
	Delete(ctx context.Context, id uint64) error
}

func (t *TransactionRepository) FindAll(ctx context.Context, req any) ([]*entity.Transaction, error) {
	ctxSpan, span := trace.StartSpan(ctx, "TransactionRepository - FindAll")
	defer span.End()

	var transactions []*entity.Transaction
	if err := t.db.Debug().WithContext(ctxSpan).Order("created_at desc").Find(&transactions).Error; err != nil {
		log.Println("ERROR: [TransactionRepository - FindAll] Internal server error:", err)
		return nil, err
	}

	return transactions, nil
}

func (t *TransactionRepository) FindByConsumerId(ctx context.Context, consumerId uint64) ([]*entity.Transaction, error) {
	ctxSpan, span := trace.StartSpan(ctx, "TransactionRepository - FindByConsumerId")
	defer span.End()

	var transactions []*entity.Transaction
	if err := t.db.Debug().WithContext(ctxSpan).Where("consumer_id = ?", consumerId).Order("created_at desc").Find(&transactions).Error; err != nil {
		log.Println("ERROR: [TransactionRepository - FindByConsumerId] Internal server error:", err)
		return nil, err
	}

	return transactions, nil
}

func (t *TransactionRepository) FindById(ctx context.Context, id uint64) (*entity.Transaction, error) {
	ctxSpan, span := trace.StartSpan(ctx, "TransactionRepository - FindById")
	defer span.End()

	var transaction entity.Transaction
	if err := t.db.Debug().WithContext(ctxSpan).Where("id = ?", id).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("WARNING: [TransactionRepository - FindById] Transaction not found for id:", id)
			return nil, status.Errorf(codes.NotFound, "Transaction not found for id: %v", id)
		}
		log.Println("ERROR: [TransactionRepository - FindById] Internal server error:", err)
		return nil, err
	}

	return &transaction, nil
}

func (t *TransactionRepository) FindByContractNumber(ctx context.Context, contractNumber string) (*entity.Transaction, error) {
	ctxSpan, span := trace.StartSpan(ctx, "TransactionRepository - FindByContractNumber")
	defer span.End()

	var transaction entity.Transaction
	if err := t.db.Debug().WithContext(ctxSpan).Where("contract_number = ?", contractNumber).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("WARNING: [TransactionRepository - FindByContractNumber] Transaction not found for contract number:", contractNumber)
			return nil, status.Errorf(codes.NotFound, "Transaction not found for contract number: %v", contractNumber)
		}
		log.Println("ERROR: [TransactionRepository - FindByContractNumber] Internal server error:", err)
		return nil, err
	}

	return &transaction, nil
}

func (t *TransactionRepository) Create(ctx context.Context, req *entity.Transaction) (*entity.Transaction, error) {
	ctxSpan, span := trace.StartSpan(ctx, "TransactionRepository - Create")
	defer span.End()

	if err := t.db.Debug().WithContext(ctxSpan).Create(req).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			log.Println("WARNING: [TransactionRepository - Create] Transaction already exists for contract number:", req.ContractNumber)
			return nil, status.Errorf(codes.AlreadyExists, "Transaction already exists for contract number: %v", req.ContractNumber)
		}
		log.Println("ERROR: [TransactionRepository - Create] Internal server error:", err)
		return nil, err
	}

	return req, nil
}

func (t *TransactionRepository) Delete(ctx context.Context, id uint64) error {
	ctxSpan, span := trace.StartSpan(ctx, "TransactionRepository - Delete")
	defer span.End()

	if err := t.db.Debug().WithContext(ctxSpan).Where("id = ?", id).Delete(&entity.Transaction{}).Error; err != nil {
		log.Println("ERROR: [TransactionRepository - Delete] Internal server error:", err)
		return err
	}

	return nil
}
