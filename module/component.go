package module

import (
	"context"
	"github.com/ganeryao/linking-go-agile/protos"
	"github.com/ganeryao/linking-go-socket/manager"
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/component"
	"github.com/topfreegames/pitaya/config"
	"github.com/topfreegames/pitaya/groups"
)

type SelfComponent interface {
	component.Component
	Group() string
	LeaveGroup(ctx context.Context, group string, uid string)
	LeaveGroupRemote(ctx context.Context, request *protos.LRequest) (*protos.LResult, error)
}

// Base implements a default module for Component.
type SelfBase struct {
	component.Base
}

func (b *SelfBase) InitGroup(conf *config.Config, group string, clearUids bool) {
	var gsi groups.GroupService
	var err error
	if conf != nil {
		gsi, err = groups.NewEtcdGroupService(conf, nil)
	} else {
		gsi = groups.NewMemoryGroupService(conf)
	}
	if err != nil {
		panic(err)
	}
	pitaya.InitGroups(gsi)
	if clearUids {
		_ = pitaya.GroupRemoveAll(context.Background(), group)
	}
	_ = pitaya.GroupCreate(context.Background(), group)
}

func (b *SelfBase) JoinGroup(ctx context.Context, IsFrontend bool, group string, uid string) error {
	logger := manager.GetLog(ctx)
	// 1、从ctx中获得session
	s := manager.GetSession(ctx)
	// 2、绑定session用户编号
	err := s.Bind(ctx, uid)
	if err != nil {
		return pitaya.Error(err, "RH-000", map[string]string{"failed": "bind"})
	}
	// 3、判断用户是否已经在组中，如果不存在再加入
	flag, err := pitaya.GroupContainsMember(ctx, group, uid)
	if err != nil {
		logger.Error("Failed to contains group member: " + err.Error())
		return err
	}
	if !flag {
		err = pitaya.GroupAddMember(ctx, group, uid)
		if err != nil {
			logger.Error("Failed to join group: " + err.Error())
			return err
		}
	}
	if IsFrontend {
		s := manager.GetSession(ctx)
		_ = s.OnClose(func() {
			logger.Error("Session Close uid : " + s.UID())
			b.LeaveGroup(ctx, group, s.UID())
		})
	}
	return nil
}

func (b *SelfBase) LeaveGroup(ctx context.Context, group string, uid string) {
	logger := manager.GetLog(ctx)
	// 1、用户从组中移除
	err := pitaya.GroupRemoveMember(ctx, group, uid)
	if err != nil {
		logger.Error("Failed to leave group member: " + err.Error())
	}
}

func (b *SelfBase) ClearGroup(ctx context.Context, group string) {
	logger := manager.GetLog(ctx)
	// 1、清空组中成员
	err := pitaya.GroupRemoveAll(ctx, group)
	if err != nil {
		logger.Error("Failed to clear group: " + err.Error())
	}
}
