package lights

import "time"

func StaticColorTempSchedule(timezone string) (ColorTempSchedulerFunc, error) {
	firstlight, err := timeOfDayToday("02:00:00", timezone)
	if err != nil {
		return nil, err
	}
	morning, err := timeOfDayToday("13:00:00", timezone)
	if err != nil {
		return nil, err
	}
	day, err := timeOfDayToday("16:00:00", timezone)
	if err != nil {
		return nil, err
	}
	lateafternoon, err := timeOfDayToday("23:00:00", timezone)
	if err != nil {
		return nil, err
	}
	evening, err := timeOfDayToday("02:00:00", timezone)
	if err != nil {
		return nil, err
	}

	return func() ColorTempSchedule {
		return ColorTempSchedule{
			ColorTemperature_FIRSTLIGHT:    *firstlight,
			ColorTemperature_MORNING:       *morning,
			ColorTemperature_DAY:           *day,
			ColorTemperature_LATEAFTERNOON: *lateafternoon,
			ColorTemperature_EVENING:       *evening,
		}
	}, nil
}

// timeOfDayToday takes a time, and replaces the date to be the same time, but today.
func timeOfDayToday(t string, timezone string) (*time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	timestamp, err := time.ParseInLocation("15:04:05", t, loc)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	d := time.Date(now.Year(), now.Month(), now.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second(), 0, loc)

	return &d, nil
}
