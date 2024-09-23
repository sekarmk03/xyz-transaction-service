package repository_test

import (
	"context"
	"regexp"
	"testing"
	"time"
	"xyz-transaction-service/modules/transaction/entity"
	"xyz-transaction-service/modules/transaction/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	return gormDB, mock, err
}

func TestFindAll(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` ORDER BY created_at desc")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "contract_number", "consumer_id", "tenor", "otr", "admin_fee", "installment", "interest", "asset_name"}).
			AddRow(1, "CN123", 1, 3, 100000, 6000, 45000, 4000, "Smartphone").
			AddRow(2, "CN124", 2, 6, 200000, 12000, 90000, 8000, "Laptop"))

	repo := repository.NewTransactionRepository(db)

	result, err := repo.FindAll(context.Background(), nil)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, uint64(1), result[0].ConsumerId)
	assert.Equal(t, "CN123", result[0].ContractNumber)
	assert.Equal(t, uint32(3), result[0].Tenor)
	assert.Equal(t, uint64(100000), result[0].Otr)
	assert.Equal(t, uint64(6000), result[0].AdminFee)
	assert.Equal(t, uint64(45000), result[0].Installment)
	assert.Equal(t, uint64(4000), result[0].Interest)
	assert.Equal(t, "Smartphone", result[0].AssetName)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestFindById(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE id = ? ORDER BY `transactions`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "contract_number", "consumer_id", "tenor", "otr", "admin_fee", "installment", "interest", "asset_name"}).
			AddRow(1, "CN124", 2, 6, 200000, 12000, 90000, 8000, "Laptop"))

	repo := repository.NewTransactionRepository(db)

	result, err := repo.FindById(context.Background(), 1)

	if result != nil {
		assert.Equal(t, uint64(1), result.Id)
		assert.Equal(t, "CN124", result.ContractNumber)
		assert.Equal(t, uint64(2), result.ConsumerId)
		assert.Equal(t, uint32(6), result.Tenor)
		assert.Equal(t, uint64(200000), result.Otr)
		assert.Equal(t, uint64(12000), result.AdminFee)
		assert.Equal(t, uint64(90000), result.Installment)
		assert.Equal(t, uint64(8000), result.Interest)
		assert.Equal(t, "Laptop", result.AssetName)
	}

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCreate(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `transactions` (`contract_number`,`consumer_id`,`tenor`,`otr`,`admin_fee`,`installment`,`interest`,`asset_name`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?,?)")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	repo := repository.NewTransactionRepository(db)

	tx := &entity.Transaction{
		ConsumerId:     3,
		ContractNumber: "CN123",
		Tenor:          12,
		Otr:            300000,
		AdminFee:       18000,
		Installment:    135000,
		Interest:       12000,
		AssetName:      "Smartwatch",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	result, err := repo.Create(context.Background(), tx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint64(3), result.ConsumerId)
	assert.Equal(t, "CN123", result.ContractNumber)
	assert.Equal(t, uint32(12), result.Tenor)
	assert.Equal(t, uint64(300000), result.Otr)
	assert.Equal(t, uint64(18000), result.AdminFee)
	assert.Equal(t, uint64(135000), result.Installment)
	assert.Equal(t, uint64(12000), result.Interest)
	assert.Equal(t, "Smartwatch", result.AssetName)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
