package queue

import "linking-go-socket/component"

type (
	Queue interface {
		Init()
		PushMsg(*component.HandlerMsg)
	}
)

var processQueue Queue

func SetProcessQueue(queue Queue) {
	processQueue = queue
}

func GetProcessQueue() Queue {
	return processQueue
}
