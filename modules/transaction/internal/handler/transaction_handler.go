package handler

import (
	"context"
	"log"
	"net/http"
	"xyz-transaction-service/common/config"
	commonErr "xyz-transaction-service/common/error"
	"xyz-transaction-service/modules/transaction/client"
	"xyz-transaction-service/modules/transaction/entity"
	"xyz-transaction-service/modules/transaction/service"
	"xyz-transaction-service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TransactionHandler struct {
	pb.UnimplementedTransactionServiceServer
	config         config.Config
	transactionSvc service.TransactionServiceUseCase
	consumerLimitSvc client.ConsumerLimitServiceClient
}

func NewTransactionHandler(config config.Config, transactionSvc service.TransactionServiceUseCase, consumerLimitSvc client.ConsumerLimitServiceClient) *TransactionHandler {
	return &TransactionHandler{
		config:         config,
		transactionSvc: transactionSvc,
		consumerLimitSvc: consumerLimitSvc,
	}
}

func (th *TransactionHandler) GetAllTransactions(ctx context.Context, req *emptypb.Empty) (*pb.TransactionListResponse, error) {
	transactionList, err := th.transactionSvc.FindAll(ctx, req)
	if err != nil {
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [TransactionHandler - GetAllTransactions] Error while find all transaction:", parseError.Message)
		return &pb.TransactionListResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	var transactions []*pb.Transaction
	for _, t := range transactionList {
		transactions = append(transactions, entity.ConvertEntityToProto(t))
	}

	return &pb.TransactionListResponse{
		Code:    uint32(http.StatusOK),
		Message: "Success get all transactions",
		Data:    transactions,
	}, nil
}

func (th *TransactionHandler) GetTransactionByContractNumber(ctx context.Context, req *pb.TransactionContractNumberRequest) (*pb.TransactionResponse, error) {
	transaction, err := th.transactionSvc.FindByContractNumber(ctx, req.ContractNumber)
	if err != nil {
		if transaction == nil {
			log.Println("WARNING: [TransactionHandler - GetTransactionByContractNumber] Transaction not found for contract number:", req.ContractNumber)
			return &pb.TransactionResponse{
				Code:    uint32(http.StatusNotFound),
				Message: "Transaction not found",
			}, status.Errorf(codes.NotFound, "Transaction not found")
		}
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [TransactionHandler - GetTransactionByContractNumber] Error while find transaction by contract number:", parseError.Message)
		return &pb.TransactionResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	return &pb.TransactionResponse{
		Code:    uint32(http.StatusOK),
		Message: "Success get transaction by contract number",
		Data:    entity.ConvertEntityToProto(transaction),
	}, nil
}

func (th *TransactionHandler) GetTransactionsByConsumerId(ctx context.Context, req *pb.TransactionConsumerIdRequest) (*pb.TransactionListResponse, error) {
	transactionList, err := th.transactionSvc.FindByConsumerId(ctx, req.ConsumerId)
	if err != nil {
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [TransactionHandler - GetTransactionsByConsumerId] Error while find transactions by consumer id:", parseError.Message)
		return &pb.TransactionListResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	var transactions []*pb.Transaction
	for _, t := range transactionList {
		transactions = append(transactions, entity.ConvertEntityToProto(t))
	}

	return &pb.TransactionListResponse{
		Code:    uint32(http.StatusOK),
		Message: "Success get transactions by consumer id",
		Data:    transactions,
	}, nil
}

func (th *TransactionHandler) CreateTransaction(ctx context.Context, req *pb.Transaction) (*pb.TransactionResponse, error) {
	// check limit available
	consumerLimit, err := th.consumerLimitSvc.GetConsumerLimitByConsumerIdAndTenor(ctx, req.ConsumerId, req.Tenor)
	if err != nil {
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [TransactionHandler - CreateTransaction] Error while get consumer limit by consumer id and tenor:", parseError.Message)
		return &pb.TransactionResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	if consumerLimit.Data.LimitAvailable < req.Otr {
		log.Println("WARNING: [TransactionHandler - CreateTransaction] Limit available not enough for consumer id:", req.ConsumerId)
		return &pb.TransactionResponse{
			Code:    uint32(http.StatusBadRequest),
			Message: "Limit available not enough",
		}, status.Errorf(codes.InvalidArgument, "Limit available not enough")
	}

	transaction, err := th.transactionSvc.Create(ctx, req.ConsumerId, req.Tenor, req.Otr, req.AdminFee, req.Installment, req.Interest, req.AssetName)
	if err != nil {
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [TransactionHandler - CreateTransaction] Error while create transaction:", parseError.Message)
		return &pb.TransactionResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	// update limit available
	_, err = th.consumerLimitSvc.UpdateAvailableLimit(ctx, req.ConsumerId, req.Tenor, req.Otr)
	if err != nil {
		parseError := commonErr.ParseError(err)
		log.Println("ERROR: [TransactionHandler - CreateTransaction] Error while update available limit:", parseError.Message)

		// should rollback transaction

		return &pb.TransactionResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	return &pb.TransactionResponse{
		Code:    uint32(http.StatusOK),
		Message: "Success create transaction",
		Data:    entity.ConvertEntityToProto(transaction),
	}, nil
}
