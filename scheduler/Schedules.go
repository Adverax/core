package scheduler

import (
	"fmt"
	"github.com/robfig/cron"
	"roketus/go-common/worker"
	"time"
)

type OnceSchedule struct {
	tick time.Time
}

func (s *OnceSchedule) NextTick(now time.Time) time.Time {
	t := s.tick
	s.tick = time.Time{}
	return t
}

func (s *OnceSchedule) String() string {
	return fmt.Sprintf("once %s", s.tick.String())
}

type PeriodicSchedule struct {
	period time.Duration
}

func (s *PeriodicSchedule) NextTick(now time.Time) time.Time {
	return now.Add(s.period)
}

func (s *PeriodicSchedule) String() string {
	return fmt.Sprintf("every %s", s.period.String())
}

func NewPeriodicSchedule(period time.Duration) Schedule {
	return &PeriodicSchedule{
		period: period,
	}
}

type CronSchedule struct {
	cron cron.Schedule
	spec string
}

func (s *CronSchedule) NextTick(now time.Time) time.Time {
	return s.cron.Next(now)
}

func (s *CronSchedule) String() string {
	return fmt.Sprintf("cron %q", s.cron)
}

func NewCronSchedule(spec string) (worker.Schedule, error) {
	s, err := cron.ParseStandard(spec)
	if err != nil {
		return nil, err
	}
	return &CronSchedule{
		cron: s,
		spec: spec,
	}, nil
}
