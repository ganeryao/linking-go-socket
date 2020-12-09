package module

import (
	"github.com/alecthomas/log4go"
	"github.com/ganeryao/linking-go-agile/common"
	"github.com/ganeryao/linking-go-agile/protos"
	"github.com/topfreegames/pitaya/util"
	"reflect"
)

type ApiProcessMode string

const (
	ApiModeMain   ApiProcessMode = "Main"
	ApiModeThread ApiProcessMode = "Thread"
	ApiModeNone   ApiProcessMode = "None"
)

type HandlerMsg struct {
	Uid     string
	ApiType ApiProcessMode
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
