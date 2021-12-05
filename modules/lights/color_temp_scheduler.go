package lights

import "time"

type ColorTempSchedulerFunc func() ColorTempSchedule

type ColorTempSchedule map[ColorTemperature]time.Time

func (c ColorTempSchedule) MostRecent() ColorTemperature {
	var mostRecent ColorTemperature

	// Find the most recent time and return the ColorTemperature
	for temp, t := range c {
		if time.Since(t) > 0 {
			if time.Since(t) < time.Since(c[mostRecent]) {
				mostRecent = temp
			}
		}
	}

	return mostRecent
}
