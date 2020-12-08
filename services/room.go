package services

import (
	"context"
	lkCommon "github.com/ganeryao/linking-go-agile/common"
	"github.com/ganeryao/linking-go-agile/protos"
	"github.com/ganeryao/linking-go-socket/app"
	"github.com/ganeryao/linking-go-socket/common"
	lkComponent "github.com/ganeryao/linking-go-socket/component"
	"github.com/ganeryao/linking-go-socket/manager"
	"github.com/ganeryao/linking-go-socket/pojo/dto"
	"github.com/spf13/viper"
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/component"
	"github.com/topfreegames/pitaya/config"
	"github.com/topfreegames/pitaya/groups"
	"github.com/topfreegames/pitaya/timer"
	"time"
)

type (
	// Room represents a component that contains a bundle of room related handler
	// like Join/Message
	Room struct {
		component.Base
		timer  *timer.Timer
		config *viper.Viper
	}

	RoomRemote struct {
		component.Base
	}
)

// NewRoom returns a new room
func NewRoom(config *viper.Viper) *Room {
	return &Room{config: config}
}

// Init runs on service initialization
func (r *Room) Init() {
	gsi, err := groups.NewEtcdGroupService(config.NewConfig(r.config), nil)
	if err != nil {
		panic(err)
	}
	pitaya.InitGroups(gsi)
	_ = pitaya.GroupCreate(context.Background(), app.DefaultGroupName)
}

// AfterInit component lifetime callback
func (r *Room) AfterInit() {
	r.timer = pitaya.NewTimer(time.Minute*5, func() {
		count, err := pitaya.GroupCountMembers(context.Background(), app.DefaultGroupName)
		println("UserCount: Time=>", time.Now().String(), "Count=>", count, "Error=>", err)
	})
}

// Join room
func (r *Room) Join(ctx context.Context, request *protos.LRequest) (*protos.LResult, error) {
	logger := manager.GetLog(ctx)
	// 1、解析加入房间请求参数
	var joinDTO = &dto.JoinDTO{}
	lkCommon.ParseJson(request.Param, joinDTO)
	// 2、从ctx中获得session
	s := manager.GetSession(ctx)
	// 3、绑定session用户编号
	logger.Info("join==============" + joinDTO.Uid)
	err := s.Bind(ctx, joinDTO.Uid)
	if err != nil {
		return nil, pitaya.Error(err, "RH-000", map[string]string{"failed": "bind"})
	}
	// 4、判断用户是否已经在组中，如果存在先删除，再加入
	flag, err := pitaya.GroupContainsMember(ctx, app.DefaultGroupName, s.UID())
	if err != nil {
		logger.Error("Failed to contains room member: " + err.Error())
		return nil, err
	}
	if flag {
		err := pitaya.GroupRemoveMember(ctx, app.DefaultGroupName, s.UID())
		if err != nil {
			logger.Error("Failed to remove room member: " + err.Error())
			return nil, err
		}
	}
	err = pitaya.GroupAddMember(ctx, app.DefaultGroupName, s.UID())
	if err != nil {
		logger.Error("Failed to join room: " + err.Error())
		return nil, err
	}
	return lkCommon.ResultOk, nil
}

// Message sync last message to all members
func (r *Room) Message(ctx context.Context, request *protos.LRequest) {
	request.Api = common.ConvertApi(request.GetApi())
	var msgDTO = &dto.MsgDTO{}
	lkCommon.ParseJson(request.Param, msgDTO)
	flag, apiType := common.ContainsApi(request.GetApi())
	if !flag {
		// 不是支持的api请求，直接抛弃
		return
	}
	s := manager.GetSession(ctx)
	// 开始处理请求
	app.GetProcessQueue().PushMsg(&lkComponent.HandlerMsg{Uid: s.UID(), ApiType: apiType, Req: request})
}
