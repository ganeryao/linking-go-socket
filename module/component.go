package module

import (
	"context"
	"github.com/ganeryao/linking-go-socket/manager"
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/component"
	"github.com/topfreegames/pitaya/config"
	"github.com/topfreegames/pitaya/groups"
)

type SelfComponent interface {
	component.Component
	Group() string
}

// Base implements a default module for Component.
type SelfBase struct {
	component.Base
}

type GroupMsg struct {
	Group string
	Uid   string
	Msg   string
}

func (b *SelfBase) InitGroup(conf *config.Config, group string) {
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
	// 初始和创建组
	pitaya.InitGroups(gsi)
	_ = pitaya.GroupCreate(context.Background(), group)
}

/**
绑定用户
*/
func (b *SelfBase) BindUser(ctx context.Context, uid string) error {
	// 1、从ctx中获得session
	s := manager.GetSession(ctx)
	// 2、绑定session用户编号
	err := s.Bind(ctx, uid)
	if err != nil {
		return pitaya.Error(err, "RH-000", map[string]string{"failed": "bind"})
	}
	return nil
}

/**
加入组
*/
func (b *SelfBase) JoinGroup(ctx context.Context, group string, uid string) bool {
	logger := manager.GetLog(ctx)
	// 1、判断用户是否已经在组中，如果不存在再加入
	flag, err := pitaya.GroupContainsMember(ctx, group, uid)
	if err != nil {
		logger.Error("Failed to contains group member: " + err.Error())
		return false
	}
	if !flag {
		err = pitaya.GroupAddMember(ctx, group, uid)
		if err != nil {
			logger.Error("Failed to join group: " + err.Error())
			return false
		}
	}
	return true
}

/**
离开组
*/
func (b *SelfBase) LeaveGroup(ctx context.Context, group string, uid string) bool {
	logger := manager.GetLog(ctx)
	// 1、用户从组中移除
	err := pitaya.GroupRemoveMember(ctx, group, uid)
	if err != nil {
		logger.Error("Failed to leave group member: " + err.Error())
		return false
	}
	return true
}

/**
清空组
*/
func (b *SelfBase) ClearGroup(ctx context.Context, group string) bool {
	logger := manager.GetLog(ctx)
	// 1、清空组中成员
	err := pitaya.GroupRemoveAll(ctx, group)
	if err != nil {
		logger.Error("Failed to clear group: " + err.Error())
		return false
	}
	return true
}
