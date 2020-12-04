package local

import (
	"github.com/alecthomas/log4go"
	"github.com/ganeryao/linking-go-socket/common"
	"github.com/ganeryao/linking-go-socket/component"
	"time"
)

var mainHandleQueue = common.GetQueue()
var threadHandleQueue = common.GetQueue()

type Queue struct {
}

func NewQueue() *Queue {
	var h = &Queue{}
	h.init()
	return h
}

func (h Queue) init() {
	// 1、1个主逻辑线程
	go h.GetNextMsg(mainHandleQueue)
	// 2、2个子逻辑线程
	for i := 0; i < 2; i++ {
		go h.GetNextMsg(threadHandleQueue)
	}
}

func (h Queue) PushMsg(handleMsg *component.HandlerMsg) {
	switch handleMsg.ApiType {
	case common.ApiModeMain:
		mainHandleQueue.Push(handleMsg)
		break
	case common.ApiModeThread:
		threadHandleQueue.Push(handleMsg)
		break
	default:
		// 不执行
		break
	}
}

func (h Queue) GetNextMsg(queue *common.Queue) {
	for {
		time.Sleep(time.Duration(2) * time.Second)
		data, _ := queue.Pop()
		if data != nil {
			var msg = data.(component.HandlerMsg)
			result, err := component.DoHandleMsg(msg)
			if err == nil {
				log4go.Info("HandleNextMsg ok=========" + msg.Uid)
			} else {
				log4go.Info("HandleNextMsg fail=========" + result.Data)
			}
		}
	}
}
