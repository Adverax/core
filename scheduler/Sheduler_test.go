package scheduler

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewPeriodicSchedule(t *testing.T) {
	interval := time.Second
	s := NewPeriodicSchedule(interval)
	now := time.Now()
	actual := s.NextTick(now).Unix()
	expected := now.Add(interval).Unix()
	require.Equal(t, expected, actual)
}

func TestScheduler(t *testing.T) {
	doneCh := make(chan struct{})
	ticks := make(map[string]int)

	s := NewScheduler(doneCh, nil)
	require.NotNil(t, s)
	s.Register(NewAction(
		"class1",
		NewPeriodicSchedule(10*time.Millisecond),
		func(ctx context.Context) error {
			ticks["class1"]++
			return nil
		},
	))

	s.Register(NewAction(
		"class2",
		NewPeriodicSchedule(30*time.Millisecond),
		func(ctx context.Context) error {
			ticks["class2"]++
			return nil
		},
	))

	s.Register(NewAction(
		"class3",
		NewPeriodicSchedule(40*time.Millisecond),
		func(ctx context.Context) error {
			ticks["class3"]++
			return nil
		},
	))
	s.Start(context.Background())

	time.Sleep(100 * time.Millisecond)
	require.Equal(t, map[string]int{
		"class1": 9,
		"class2": 3,
		"class3": 2,
	}, ticks)
	close(doneCh)
}
