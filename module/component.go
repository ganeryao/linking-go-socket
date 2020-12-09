package module

import "github.com/topfreegames/pitaya/component"

type SelfComponent interface {
	component.Component
	Group() string
}

// Base implements a default module for Component.
type SelfBase struct {
	component.Base
}
