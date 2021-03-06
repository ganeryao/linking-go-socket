package manager

import (
	"context"
	"fmt"
	"github.com/ganeryao/linking-go-agile/common"
	"github.com/ganeryao/linking-go-agile/protos"
	"github.com/topfreegames/pitaya"
	"github.com/topfreegames/pitaya/session"
)

// SessionData struct
type SessionData struct {
	Data map[string]interface{}
}

type UserSession struct {

}

// GetSession gets the session data
func GetSession(ctx context.Context) *session.Session {
	return pitaya.GetSessionFromCtx(ctx)
}

// GetSessionData gets the session data
func (c *UserSession) GetSessionData(ctx context.Context) (*SessionData, error) {
	s := pitaya.GetSessionFromCtx(ctx)
	res := &SessionData{
		Data: s.GetData(),
	}
	return res, nil
}

// SetSessionData sets the session data
func (c *UserSession) SetSessionData(ctx context.Context, data *SessionData) (*protos.LResult, error) {
	s := pitaya.GetSessionFromCtx(ctx)
	err := s.SetData(data.Data)
	if err != nil {
		return nil, pitaya.Error(err, "CN-000", map[string]string{"failed": "set data"})
	}
	err = s.PushToFront(ctx)
	if err != nil {
		return nil, err
	}
	return common.ResultOk, nil
}

// NotifySessionData sets the session data
func (c *UserSession) NotifySessionData(ctx context.Context, data *SessionData) {
	s := pitaya.GetSessionFromCtx(ctx)
	err := s.SetData(data.Data)
	if err != nil {
		fmt.Println("got error on notify", err)
	}
}