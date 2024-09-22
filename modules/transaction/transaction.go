package transaction

import (
	"xyz-transaction-service/common/config"
	"xyz-transaction-service/modules/transaction/internal/builder"
	"xyz-transaction-service/pb"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func InitGrpc(server *grpc.Server, cfg config.Config, db *gorm.DB, grpcConn *grpc.ClientConn) {
	transaction := builder.BuildTransactionHandler(cfg, db, grpcConn)
	pb.RegisterTransactionServiceServer(server, transaction)
}
