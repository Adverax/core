package scheduler

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

type Logger interface {
	LogError(ctx context.Context, err error)
}

type Scheduler struct {
	mx        sync.Mutex
	logger    Logger
	events    events
	wakeup    chan struct{}
	appDoneCh chan struct{}
}

func (s *Scheduler) Register(e Event) {
	s.update(&event{Event: e}, true)
}

func (s *Scheduler) Start(ctx context.Context) {
	go s.work(ctx)
}

func (s *Scheduler) work(ctx context.Context) {
	for {
		event := s.capture()

		select {
		case <-s.appDoneCh:
			return
		case <-s.wakeup:
			s.update(event, false)
			continue
		case <-time.After(getSleepTime(event)):
			if event == nil {
				continue
			}
		}

		s.handle(ctx, event)
	}
}

func (s *Scheduler) handle(ctx context.Context, event *event) {
	err := event.Trigger(ctx)
	if err != nil {
		s.logger.LogError(ctx, err)
	}

	s.update(event, true)
}

func (s *Scheduler) display() {
	fmt.Println(fmt.Sprintf("Scheduler events: %s", time.Now().String()))
	for _, event := range s.events {
		fmt.Println(event.String())
	}
}

func (s *Scheduler) capture() *event {
	s.mx.Lock()
	defer s.mx.Unlock()

	if len(s.events) == 0 {
		return nil
	}

	e := s.events[0]
	s.events = s.events[1:]

	return e
}

func (s *Scheduler) update(event *event, recalc bool) {
	if event == nil {
		return
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	s.refresh(event, recalc)
}

func (s *Scheduler) refresh(event *event, recalc bool) {
	if recalc {
		s.recalc(event)
	}
	if event.nextTick.IsZero() {
		return
	}
	s.events = insertIntoEvents(s.events, event)
}

func (s *Scheduler) recalc(event *event) {
	event.nextTick = event.NextTick(time.Now())
	event.nextTimestamp = event.nextTick.UnixNano()
}

func NewScheduler(
	appDoneCh chan struct{},
	logger Logger,
) *Scheduler {
	return &Scheduler{
		wakeup:    make(chan struct{}, 100),
		appDoneCh: appDoneCh,
		logger:    logger,
	}
}

func getSleepTime(event *event) time.Duration {
	if event == nil {
		return time.Minute
	}

	now := time.Now()
	if event.nextTimestamp < now.UnixNano() {
		return time.Millisecond
	}

	return event.nextTick.Sub(now)
}

func insertIntoEvents(events events, event *event) events {
	index := sort.Search(
		len(events),
		func(i int) bool {
			return events[i].nextTimestamp >= event.nextTimestamp
		},
	)
	events = append(events, nil)
	copy(events[index+1:], events[index:])
	events[index] = event
	return events
}
