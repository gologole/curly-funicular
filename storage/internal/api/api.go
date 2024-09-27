package api

import "cmd/main.go/configs"

type ELKAPI interface {
	SendData(data []byte) ([]byte, error)
}

type ELKAPIImpl struct {
	cfgs *configs.Config
}
