package scheduler

import (
	"context"
	"fmt"
	"time"
)

type Schedule interface {
	NextTick(now time.Time) time.Time
	String() string
}

type Event interface {
	Schedule
	Name() string
	Trigger(ctx context.Context) error
	String() string
}

type Action struct {
	schedule Schedule
	name     string
	handler  func(ctx context.Context) error
	active   bool
	async    bool
}

func (a *Action) String() string {
	return fmt.Sprintf("Action %s with schedule %s", a.name, a.schedule.String())
}

func (a *Action) Name() string {
	return a.name
}

func (a *Action) NextTick(now time.Time) time.Time {
	return a.schedule.NextTick(now)
}

func NewAction(name string, schedule Schedule, handler func(ctx context.Context) error) *Action {
	return &Action{
		schedule: schedule,
		name:     name,
		handler:  handler,
	}
}

type AsyncAction struct {
	schedule Schedule
	name     string
	handler  func(ctx context.Context) error
	active   bool
}

func (a *AsyncAction) String() string {
	return fmt.Sprintf("Action %s with schedule %s", a.name, a.schedule.String())
}

func (a *AsyncAction) Name() string {
	return a.name
}

func (a *AsyncAction) NextTick(now time.Time) time.Time {
	return a.schedule.NextTick(now)
}

func (a *AsyncAction) Trigger(ctx context.Context) error {
	if a.active {
		return nil
	}

	a.active = true
	go func() {
		defer func() {
			a.active = false
		}()

		err := a.handler(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()

	return nil
}

func NewAsyncAction(name string, schedule Schedule, handler func(ctx context.Context) error) *AsyncAction {
	return &AsyncAction{
		schedule: schedule,
		name:     name,
		handler:  handler,
	}
}

type event struct {
	Event
	nextTimestamp int64
	nextTick      time.Time
}

func (e *event) String() string {
	return fmt.Sprintf("%s will ticked at %d", e.Event.String(), e.nextTimestamp)
}

type events []*event

func (es events) Len() int {
	return len(es)
}

func (es events) Less(i, j int) bool {
	return es[i].nextTimestamp < es[j].nextTimestamp
}

func (es events) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}
