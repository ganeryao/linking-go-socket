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

func (b *SelfBase) Group() string {
	return ""
}

func (b *SelfBase) InitGroup(conf *config.Config, clearUids bool) {
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
		_ = pitaya.GroupRemoveAll(context.Background(), b.Group())
	}
	_ = pitaya.GroupCreate(context.Background(), b.Group())
}

func (b *SelfBase) JoinGroup(ctx context.Context, uid string) error {
	logger := manager.GetLog(ctx)
	// 1、从ctx中获得session
	s := manager.GetSession(ctx)
	// 2、绑定session用户编号
	logger.Info("join group ==============" + uid)
	err := s.Bind(ctx, uid)
	if err != nil {
		return pitaya.Error(err, "RH-000", map[string]string{"failed": "bind"})
	}
	// 3、判断用户是否已经在组中，如果不存在再加入
	flag, err := pitaya.GroupContainsMember(ctx, b.Group(), uid)
	if err != nil {
		logger.Error("Failed to contains room member: " + err.Error())
		return err
	}
	if !flag {
		err = pitaya.GroupAddMember(ctx, b.Group(), uid)
		if err != nil {
			logger.Error("Failed to join room: " + err.Error())
			return err
		}
	}
	return nil
}
