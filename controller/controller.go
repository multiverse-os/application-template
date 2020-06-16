package controller

type Action struct {
}

type Controller interface {
	BeforeHooks() *Action
	AfterHooks() *Action
	Actions() *Action
	Watch() err
}
