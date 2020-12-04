package queue

import "github.com/ganeryao/linking-go-socket/component"

type (
	LkQueue interface {
		Init()
		PushMsg(*component.HandlerMsg)
	}
)
