package runtime

import (
	event "github.com/swift9/ares-event"
)

type STATUS string

const (
	StatusUp   STATUS = "UP"
	StatusDown STATUS = "DOWN"
)

type Status struct {
	Status  STATUS                 `json:"status"`
	Details map[string]interface{} `json:"details"`
}

func NewStatusUp() Status {
	return Status{
		Status: StatusUp,
	}
}

func NewStatusDown() Status {
	return Status{
		Status: StatusDown,
	}
}

type IRuntime interface {
	Init()
	Start() int
	Stop() int
	Idle() int
	Health() Status
	event.IEmitter
}
