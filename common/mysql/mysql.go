package mysql

import (
	"fmt"
	"xyz-transaction-service/common/config"
)

func NewPool(cfg *config.MySQL) (string, error) {
	connCfg := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	return connCfg, nil
}
