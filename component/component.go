package component

type SelfComponent interface {
	Init()
}

// Base implements a default component for Component.
type Base struct{}

// Init was called to initialize the component.
func (c *Base) Init() {}
