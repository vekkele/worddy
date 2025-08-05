package domain

import (
	"slices"
	"time"

	"github.com/vekkele/worddy/internal/utils"
)

type ReviewsAtTime struct {
	Count int
	Time  time.Time
}

type HourForecast struct {
	Time            time.Time
	Count           int
	CumulativeCount int
}

type DayForecast struct {
	Weekday    time.Weekday
	TotalCount int
	Hourly     []HourForecast
}

func (d *DayForecast) Empty() bool {
	return len(d.Hourly) < 1 || d.TotalCount == 0
}

func (d *DayForecast) CumulativeCount() int {
	if d.Empty() {
		return 0
	}

	lastHour := d.Hourly[len(d.Hourly)-1]
	return lastHour.CumulativeCount
}

func NewForecast(reviews []ReviewsAtTime, loc *time.Location) []DayForecast {
	now := time.Now().In(loc)
	maxDate := utils.StartOfDay(now).AddDate(0, 0, 7)
	daysMap := prepareDaysMap(reviews, maxDate)

	firstWeekday := now.Weekday()
	var dayForecasts []DayForecast
	var cumulativeCount int

	for i := range 7 {
		weekday := (firstWeekday + time.Weekday(i)) % 7

		hoursMap := daysMap[weekday]
		if hoursMap == nil {
			dayForecasts = append(dayForecasts, DayForecast{
				Weekday: weekday,
			})
			continue
		}

		hours := make([]time.Time, 0, len(hoursMap))
		for hour := range hoursMap {
			hours = append(hours, hour)
		}
		slices.SortFunc(hours, func(a, b time.Time) int {
			return a.Compare(b)
		})

		hourly := make([]HourForecast, 0, len(hours))
		total := 0

		for _, hour := range hours {
			count := hoursMap[hour]
			cumulativeCount += count
			total += count

			hourly = append(hourly, HourForecast{
				Time:            hour,
				Count:           count,
				CumulativeCount: cumulativeCount,
			})
		}

		dayForecasts = append(dayForecasts, DayForecast{
			Weekday:    weekday,
			TotalCount: total,
			Hourly:     hourly,
		})
	}

	return dayForecasts
}

func prepareDaysMap(reviews []ReviewsAtTime, maxDate time.Time) map[time.Weekday]map[time.Time]int {
	daysMap := make(map[time.Weekday]map[time.Time]int)

	for _, r := range reviews {
		if !r.Time.Before(maxDate) {
			continue
		}

		weekday := r.Time.Weekday()

		if _, ok := daysMap[weekday]; !ok {
			daysMap[weekday] = make(map[time.Time]int)
		}

		daysMap[weekday][r.Time] += r.Count
	}

	return daysMap
}
