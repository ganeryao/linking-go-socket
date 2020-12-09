package common

import (
	lkCommon "github.com/ganeryao/linking-go-agile/common"
	lkError "github.com/ganeryao/linking-go-agile/errors"
	"github.com/ganeryao/linking-go-agile/protos"
	"github.com/ganeryao/linking-go-socket/constants"
	"github.com/ganeryao/linking-go-socket/linking"
	"github.com/ganeryao/linking-go-socket/module"
	"strings"
)

func ConvertApi(api string) string {
	a := strings.Split(api, ".")
	num := len(a)
	if num >= 3 {
		return a[len(a)-3] + "." + a[len(a)-2] + "." + a[len(a)-1]
	} else if num >= 2 {
		return a[len(a)-2] + "." + a[len(a)-1]
	} else {
		panic("ConvertApi api len error" + api)
	}
}

func ConvertHandlerMsg(request *protos.LRequest, uid string, data interface{}) (*module.HandlerMsg, *lkError.Error) {
	request.Api = ConvertApi(request.GetApi())
	lkCommon.ParseJson(request.Param, data)
	flag, apiType := linking.ContainsApi(request.GetApi())
	if !flag {
		// 不是支持的api请求，直接抛弃
		return nil, lkError.NewError(constants.ErrUnsupportedRequest, lkError.ErrBadRequestCode)
	}
	return &module.HandlerMsg{Uid: uid, ApiType: apiType, Api: request.Api, Msg: data}, nil
}
