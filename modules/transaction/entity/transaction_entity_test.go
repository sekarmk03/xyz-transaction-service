package entity_test

import (
	"testing"

	"xyz-transaction-service/modules/transaction/entity"
)

func TestNewTransactionEntity(t *testing.T) {
	t.Log("TestNewTransactionEntity")

	contractNumber := "CNTR-20240922-00000001-1b2b3c4d"
	consumerId := uint64(1)
	tenor := uint32(6)
	otr := uint64(50000)
	adminFee := uint64(10000)
	installment := uint64(10000)
	interest := uint64(5000)
	assetName := "Toyota Avanza 2022"
	e := entity.NewTransactionEntity(contractNumber, consumerId, tenor, otr, adminFee, installment, interest, assetName)
	if e == nil {
		t.Error("NewTransactionEntity() returned nil")
	} else {
		if e.ContractNumber != contractNumber {
			t.Error("NewTransactionEntity() returned incorrect To")
		}
		if e.ConsumerId != consumerId {
			t.Error("NewTransactionEntity() returned incorrect ConsumerId")
		}
		if e.Tenor != tenor {
			t.Error("NewTransactionEntity() returned incorrect Tenor")
		}
		if e.Otr != otr {
			t.Error("NewTransactionEntity() returned incorrect Otr")
		}
		if e.AdminFee != adminFee {
			t.Error("NewTransactionEntity() returned incorrect AdminFee")
		}
		if e.Installment != installment {
			t.Error("NewTransactionEntity() returned incorrect Installment")
		}
		if e.Interest != interest {
			t.Error("NewTransactionEntity() returned incorrect Interest")
		}
		if e.AssetName != assetName {
			t.Error("NewTransactionEntity() returned incorrect AssetName")
		}
	}
}
