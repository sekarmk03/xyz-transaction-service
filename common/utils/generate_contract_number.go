package utils

import (
	"fmt"
	"time"
	"github.com/google/uuid"
)

func GenerateContractNumber(consumerId uint64) string {
	prefix := "CNTR"
	currentDate := time.Now().Format("20060102")
	formattedConsumerID := fmt.Sprintf("%08d", consumerId)
	uniqueID := uuid.New().String()
	contractNumber := fmt.Sprintf("%s-%s-%s-%s", prefix, currentDate, formattedConsumerID, uniqueID[:8])
	return contractNumber
}
