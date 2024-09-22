package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func GenerateContractNumber(consumerId uint64) string {
	prefix := "CNTR"
	currentDate := time.Now().Format("20060102")
	formattedConsumerID := fmt.Sprintf("%08d", consumerId)
	uniqueID := strconv.Itoa(rand.Intn(1000))
	contractNumber := fmt.Sprintf("%s-%s-%s-%s", prefix, currentDate, formattedConsumerID, uniqueID)
	return contractNumber
}
