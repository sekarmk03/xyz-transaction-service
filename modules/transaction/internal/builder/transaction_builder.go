package builder

import (
	"xyz-transaction-service/common/config"
	"xyz-transaction-service/modules/transaction/internal/handler"
	"xyz-transaction-service/modules/transaction/internal/repository"
	"xyz-transaction-service/modules/transaction/service"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func BuildTransactionHandler(cfg config.Config, db *gorm.DB, grpcConn *grpc.ClientConn) *handler.TransactionHandler {
	transactionRepository := repository.NewTransactionRepository(db)
	transactionSvc := service.NewTransactionService(cfg, transactionRepository)

	return handler.NewTransactionHandler(cfg, transactionSvc)
}
