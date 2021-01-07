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
