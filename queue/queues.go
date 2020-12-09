package queue

import "github.com/ganeryao/linking-go-socket/module"

type (
	LkQueue interface {
		Init()
		PushMsg(module.HandlerMsg)
	}
)
