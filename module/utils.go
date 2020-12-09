package module

import (
	"github.com/alecthomas/log4go"
	"github.com/ganeryao/linking-go-agile/common"
	"github.com/ganeryao/linking-go-agile/protos"
	"github.com/ganeryao/linking-go-socket/linking"
	"github.com/topfreegames/pitaya/util"
	"reflect"
)

type HandlerMsg struct {
	Uid     string
	ApiType linking.ApiProcessMode
	Api     string
	Msg     interface{}
}

func DoHandleMsg(msg HandlerMsg) (*protos.LResult, error) {
	handler := handlers[msg.Api]
	args := []reflect.Value{handler.Receiver, reflect.ValueOf(msg)}
	ret, err := util.Pcall(handler.Method, args)
	if err != nil {
		_ = log4go.Error("=============" + err.Error())
		return common.OfResultFail("123", err.Error()), err
	}
	if ret != nil {
		return ret.(*protos.LResult), nil
	}
	return common.ResultOk, nil
}
