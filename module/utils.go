package module

import (
	"context"
	"github.com/alecthomas/log4go"
	"github.com/ganeryao/linking-go-agile/common"
	"github.com/ganeryao/linking-go-agile/errors"
	"github.com/ganeryao/linking-go-agile/protos"
	"github.com/topfreegames/pitaya"
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
		_ = log4go.Error("DoHandleMsg=============" + err.Error())
		return common.OfResultFail(msg.Api, errors.ErrInternalCode, err.Error()), err
	}
	if ret != nil {
		return ret.(*protos.LResult), nil
	}
	return common.ResultOk, nil
}

// SendRPC sends rpc
func SendRPC(ctx context.Context, serverId string, route string, request *protos.LRequest) {
	logger := pitaya.GetDefaultLoggerFromCtx(ctx)
	ret := &protos.LResult{}
	err := pitaya.RPCTo(ctx, serverId, route, ret, request)
	if err != nil {
		logger.Errorf("Failed to execute RPCTo %s - %s", route, err.Error())
	}
}
