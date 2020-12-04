package queue

import "github.com/ganeryao/linking-go-socket/component"

type (
	LkQueue interface {
		Init()
		PushMsg(*component.HandlerMsg)
	}
)

var processQueue LkQueue

func SetProcessQueue(queue LkQueue) {
	processQueue = queue
}

func GetProcessQueue() LkQueue {
	return processQueue
}
