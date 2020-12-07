package socket

import (
	"context"
	"fmt"
	"github.com/ganeryao/linking-go-socket/app"
	"github.com/ganeryao/linking-go-socket/services"
	"github.com/spf13/viper"
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/acceptor"
	"github.com/topfreegames/pitaya/cluster"
	"github.com/topfreegames/pitaya/component"
	"github.com/topfreegames/pitaya/route"
	"strings"
)

func ConfigureBackend(config *viper.Viper) {
	room := services.NewRoom(config)
	pitaya.Register(room,
		component.WithName(app.DefaultGroupName),
		component.WithNameFunc(strings.ToLower),
	)

	pitaya.RegisterRemote(room,
		component.WithName(app.DefaultGroupName),
		component.WithNameFunc(strings.ToLower),
	)
}

func ConfigureFrontend(port int, dictionary ...string) {
	pitaya.Register(&services.Connector{},
		component.WithName(app.DefaultFrontend),
		component.WithNameFunc(strings.ToLower),
	)
	pitaya.RegisterRemote(&services.ConnectorRemote{},
		component.WithName(app.DefaultFrontend+"_remote"),
	)
	num := len(app.DefaultRoutes)
	for i := 0; i < num; i++ {
		AddRoute(app.DefaultRoutes[i])
	}
	var dict = make(map[string]uint16, 0)
	num = len(dictionary)
	if num > 0 {
		for i := 0; i < num; i++ {
			dict[dictionary[i]] = uint16(i)
		}
	}
	err := pitaya.SetDictionary(dict)

	if err != nil {
		fmt.Printf("error setting route dictionary %s\n", err.Error())
		panic(err)
	}
	wsPort := fmt.Sprintf(":%d", port)
	tcp := acceptor.NewWSAcceptor(wsPort)
	pitaya.AddAcceptor(tcp)
}

/**
添加服务路由
*/
func AddRoute(serverType string) {
	err := pitaya.AddRoute(serverType, func(
		ctx context.Context,
		route *route.Route,
		payload []byte,
		servers map[string]*cluster.Server,
	) (*cluster.Server, error) {
		for k := range servers {
			return servers[k], nil
		}
		return nil, nil
	})
	if err != nil {
		panic(err)
	}
}
