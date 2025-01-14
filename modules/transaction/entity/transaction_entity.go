package entity

import (
	"time"
	"xyz-transaction-service/pb"
)

const (
	TransactionTableName = "transactions"
)

type Transaction struct {
	Id             uint64    `json:"id"`
	ContractNumber string    `json:"contract_number"`
	ConsumerId     uint64    `json:"consumer_id"`
	Tenor          uint32    `json:"tenor"`
	Otr            uint64    `json:"otr"`
	AdminFee       uint64    `json:"admin_fee"`
	Installment    uint64    `json:"installment"`
	Interest       uint64    `json:"interest"`
	AssetName      string    `json:"asset_name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func NewTransactionEntity(contractNumber string, consumerId uint64, tenor uint32, otr uint64, adminFee uint64, installment uint64, interest uint64, assetName string) *Transaction {
	return &Transaction{
		ContractNumber: contractNumber,
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
}

func (t *Transaction) TableName() string {
	return TransactionTableName
}

func ConvertEntityToProto(t *Transaction) *pb.Transaction {
	return &pb.Transaction{
		Id:             t.Id,
		ContractNumber: t.ContractNumber,
		ConsumerId:     t.ConsumerId,
		Tenor:          t.Tenor,
		Otr:            t.Otr,
		AdminFee:       t.AdminFee,
		Installment:    t.Installment,
		Interest:       t.Interest,
		AssetName:      t.AssetName,
		CreatedAt:      t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      t.UpdatedAt.Format(time.RFC3339),
	}
}
