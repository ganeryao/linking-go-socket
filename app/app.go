package app

import (
	"github.com/ganeryao/linking-go-socket/queue"
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/logger"
)

var (
	DefaultFrontend       = "connector"
	DefaultGroupName      = "room"
	DefaultOnMessageRoute = "onMessage"
	DefaultRoutes         = []string{"chat", "room"}
	processQueue          queue.LkQueue
)

func SetProcessQueue(queue queue.LkQueue) {
	processQueue = queue
}

func GetProcessQueue() queue.LkQueue {
	return processQueue
}

func SendUserMsg(uid string, v interface{}) error {
	return SendUsersMsg([]string{uid}, v)
}

func SendUsersMsg(uids []string, v interface{}) error {
	errUids, err := pitaya.SendPushToUsers(DefaultOnMessageRoute, v, uids, DefaultFrontend)
	if err != nil {
		logger.Log.Errorf("SendUserMsg error, UID=%d, Error=%s", errUids, err.Error())
		return err
	}
	return nil
}
