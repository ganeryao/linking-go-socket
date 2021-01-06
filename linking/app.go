package linking

import (
	"encoding/json"
	"github.com/alecthomas/log4go"
	"github.com/ganeryao/linking-go-socket/module"
	"github.com/ganeryao/linking-go-socket/queue"
	"github.com/spf13/viper"
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/logger"
)

// App is the base linking struct
type App struct {
	isFrontend  bool
	frontend    string
	config      *viper.Viper
	configured  bool
	debug       bool
	routes      []string
	clientRoute string
	queue       queue.LkQueue
}

var (
	app = &App{
		frontend:    "connector",
		clientRoute: "onMessage",
		debug:       false,
		configured:  false,
	}
)

// Configure configures the linking
func Configure(
	isFrontend bool,
	frontend string,
	clientRoute string,
	routes []string,
	queue queue.LkQueue,
	config *viper.Viper,
	debug bool,
) {
	if app.configured {
		logger.Log.Warn("lk socket configured twice!")
	}
	app.isFrontend = isFrontend
	app.frontend = frontend
	app.clientRoute = clientRoute
	app.routes = routes
	app.config = config
	app.queue = queue
	app.configured = true
	app.debug = debug
}

func IsFrontend() bool {
	return app.isFrontend
}

func GetFrontend() string {
	return app.frontend
}

func GetQueue() queue.LkQueue {
	return app.queue
}

func GetClientRoute() string {
	return app.clientRoute
}

func GetRoutes() []string {
	return app.routes
}

func IsDebug() bool {
	return app.debug
}

func IsQueueProcess() bool {
	return app.queue != nil
}

func HandleMsg(msg *module.HandlerMsg) {
	if IsQueueProcess() {
		GetQueue().PushMsg(*msg)
	} else {
		_, err := module.DoHandleMsg(*msg)
		if err != nil {
			b, _ := json.Marshal(msg)
			_ = log4go.Error("HandleMsg error=========msg:"+string(b), err)
		}
	}
}

func SendUserMsg(uid string, v interface{}) error {
	return SendUsersMsg([]string{uid}, v)
}

func SendUsersMsg(uids []string, v interface{}) error {
	errUids, err := pitaya.SendPushToUsers(app.clientRoute, v, uids, app.frontend)
	if err != nil {
		logger.Log.Errorf("SendUserMsg error, UID=%d, Error=%s", errUids, err.Error())
		return err
	}
	return nil
}
