package client

import (
	"context"
	"xyz-transaction-service/pb"
	"xyz-transaction-service/server"
)

type ConsumerLimitServiceClient struct {
	Client pb.ConsumerLimitServiceClient
}

func BuildConsumerLimitServiceClient(url string) ConsumerLimitServiceClient {
	cc := server.InitGRPCConn(url, false, "")

	c := ConsumerLimitServiceClient{
		Client: pb.NewConsumerLimitServiceClient(cc),
	}

	return c
}

func (cla *ConsumerLimitServiceClient) GetConsumerLimitByConsumerIdAndTenor(ctx context.Context, consumerId uint64, tenor uint32) (*pb.ConsumerLimitResponse, error) {
	req := &pb.ConsumerIdAndTenorRequest{
		ConsumerId: consumerId,
		Tenor:      tenor,
	}

	return cla.Client.GetConsumerLimitByConsumerIdAndTenor(ctx, req)
}

func (cla *ConsumerLimitServiceClient) UpdateAvailableLimit(ctx context.Context, consumerId uint64, tenor uint32, amountTransaction uint64) (*pb.ConsumerLimitResponse, error) {
	req := &pb.UpdateAvailableLimitRequest{
		ConsumerId:        consumerId,
		Tenor:             tenor,
		AmountTransaction: amountTransaction,
	}

	return cla.Client.UpdateAvailableLimit(ctx, req)
}
