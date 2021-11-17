package events

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// Scheduler holds timeSlice objects and provides an methods to update them..
type Scheduler struct {
	logger    log.Logger
	timeSlice *TimeSlice
}

// TimeSlice is an association between a specific time, and the names of the events that should fire at that time.
type TimeSlice map[time.Time][]string

// NewScheduler returns a new Scheduler.
func NewScheduler(logger log.Logger) *Scheduler {
	return &Scheduler{
		timeSlice: &TimeSlice{},
		logger:    logger,
	}
}

// All returns all current timeSlice objects.
func (s *Scheduler) All() TimeSlice {
	return *s.timeSlice
}

func (s *Scheduler) Report() {
	fields := []interface{}{"msg", "events"}

	for k, v := range *s.timeSlice {
		fields = append(fields, k.Format(time.RFC3339))
		fields = append(fields, strings.Join(v, ","))
	}

	_ = level.Info(s.logger).Log(fields...)
}

// Next determines the next occurring event in the series.
func (s *Scheduler) Next() *time.Time {
	times := s.ordered()

	if len(times) > 0 {
		return &times[0]
	}

	return nil
}

func (s *Scheduler) ordered() []time.Time {
	keys := []time.Time{}

	for k := range *s.timeSlice {
		keys = append(keys, k)
	}

	if len(keys) == 0 {
		return nil
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	return keys
}

// TimesForName returns all timeSlices for a given event name.
func (s *Scheduler) TimesForName(n string) []time.Time {
	var times []time.Time

	for t, names := range *s.timeSlice {
		for _, name := range names {
			if n == name {
				times = append(times, t)
			}
		}
	}

	return times
}

// NamesForTime returns all event names that are scheduled for a given timeSlice.
func (s *Scheduler) NamesForTime(t time.Time) []string {
	return (*s.timeSlice)[t]
}

// WaitForNext is a blocking function that waits for the next available time to
// arrive before returning the names to the caller.
func (s *Scheduler) WaitForNext() []string {
	next := s.Next()

	if next == nil {
		return []string{}
	}

	// Send past events under 30 seconds old.
	if time.Since(*next) > time.Duration(30)*time.Second {
		_ = level.Warn(s.logger).Log("msg", "sending past event",
			"next", next,
			"since", time.Since(*next),
		)
		return s.NamesForTime(*next)
	}

	_ = level.Info(s.logger).Log("msg", "scheduler waiting",
		"next", time.Until(*next),
		"names", strings.Join(s.NamesForTime(*next), ","),
	)

	ti := time.NewTimer(time.Until(*next))
	<-ti.C

	return s.NamesForTime(*next)
}

// Step deletes the next timeSlice.  This is determined to be the timeSlice
// that has just run.  The expectation is that Step() is called once the
// events have completed firing to advance to the next position in time.
func (s *Scheduler) Step() {
	next := s.Next()

	if next != nil {
		delete(*s.timeSlice, *s.Next())
	}
}

// Set appends the name given to the time slot given.
func (s *Scheduler) Set(t time.Time, name string) error {

	if name == "" {
		return fmt.Errorf("unable to schedule empty name at time %s", t)
	}

	if time.Until(t) < 0 {
		if time.Since(t) > 5*time.Second {
			return fmt.Errorf("not scheduling past event %s for %s, %s", name, t, time.Until(t))
		}
	}

	if _, ok := (*s.timeSlice)[t]; !ok {
		(*s.timeSlice)[t] = make([]string, 0)
	}

	_ = level.Debug(s.logger).Log("msg", "scheduling event",
		"name", name,
		"time", t,
	)

	timeHasName := func(names []string) bool {
		for _, n := range names {
			if n == name {
				return true
			}
		}

		return false
	}((*s.timeSlice)[t])

	if !timeHasName {
		(*s.timeSlice)[t] = append((*s.timeSlice)[t], name)
	}

	return nil
}
