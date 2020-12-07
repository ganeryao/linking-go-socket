package app

import "github.com/ganeryao/linking-go-socket/queue"

var processQueue queue.LkQueue

func SetProcessQueue(queue queue.LkQueue) {
	processQueue = queue
}

func GetProcessQueue() queue.LkQueue {
	return processQueue
}
