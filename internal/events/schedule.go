package events

import (
	"sort"
	"time"

	"github.com/apex/log"
)

// Scheduler
type Scheduler struct {
	timeSlice *timeSlice
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		timeSlice: &timeSlice{},
	}
}

func (s *Scheduler) All() timeSlice {
	return *s.timeSlice
}

func (s *Scheduler) Next() time.Time {

	keys := []time.Time{}

	for k, _ := range *s.timeSlice {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	return keys[0]
}

func (s *Scheduler) NamesForTime(t time.Time) []string {
	return (*s.timeSlice)[t]
}

func (s *Scheduler) WaitForNext() []string {
	next := s.Next()

	if time.Now().After(next) {
		log.Infof("sending past event: %s", next)
		return s.NamesForTime(next)
	}

	log.Infof("waiting until: %s", next)
	ti := time.NewTimer(time.Until(next))
	<-ti.C

	return s.NamesForTime(next)
}

func (s *Scheduler) Step() {
	delete(*s.timeSlice, s.Next())
}

func (s *Scheduler) Set(t time.Time, name string) {

	if _, ok := (*s.timeSlice)[t]; !ok {
		(*s.timeSlice)[t] = make([]string, 1)
	}

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

}

type timeSlice map[time.Time][]string
