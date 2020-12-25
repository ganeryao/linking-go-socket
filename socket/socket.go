package socket

import (
	"context"
	"fmt"
	"github.com/ganeryao/linking-go-socket/linking"
	"github.com/ganeryao/linking-go-socket/module"
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/acceptor"
	"github.com/topfreegames/pitaya/cluster"
	"github.com/topfreegames/pitaya/component"
	"github.com/topfreegames/pitaya/route"
	"net/http"
	"strings"
)

func ConfigureBackend(c module.SelfComponent) {
	pitaya.Register(c,
		component.WithName(c.Group()),
		component.WithNameFunc(strings.ToLower),
	)
	pitaya.RegisterRemote(c,
		component.WithName(c.Group()),
		component.WithNameFunc(strings.ToLower),
	)
}

func ConfigureFrontend(c module.SelfComponent, remote module.SelfComponent, port int, httpPort int, dictionary ...string) {
	pitaya.Register(c,
		component.WithName(c.Group()),
		component.WithNameFunc(strings.ToLower),
	)
	pitaya.RegisterRemote(remote,
		component.WithName(remote.Group()),
	)
	num := len(linking.GetRoutes())
	for i := 0; i < num; i++ {
		AddRoute(linking.GetRoutes()[i])
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
	hPort := fmt.Sprintf(":%d", httpPort)
	go http.ListenAndServe(hPort, nil)
}

func ConfigureStandalone(port int, httpPort int) {
	wsPort := fmt.Sprintf(":%d", port)
	tcp := acceptor.NewWSAcceptor(wsPort)
	pitaya.AddAcceptor(tcp)
	wsHttpPort := fmt.Sprintf(":%d", httpPort)
	go http.ListenAndServe(wsHttpPort, nil)
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
