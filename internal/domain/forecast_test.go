package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vekkele/worddy/internal/utils"
)

func TestNewForecast(t *testing.T) {
	loc := time.UTC
	now := time.Now().In(loc)
	todayStart := utils.StartOfDay(now)

	tests := []struct {
		Name   string
		Input  []ReviewsAtTime
		Result []DayForecast
	}{
		{
			Name:  "empty",
			Input: []ReviewsAtTime{},
			Result: []DayForecast{
				{Weekday: now.Weekday()},
				{Weekday: now.AddDate(0, 0, 1).Weekday()},
				{Weekday: now.AddDate(0, 0, 2).Weekday()},
				{Weekday: now.AddDate(0, 0, 3).Weekday()},
				{Weekday: now.AddDate(0, 0, 4).Weekday()},
				{Weekday: now.AddDate(0, 0, 5).Weekday()},
				{Weekday: now.AddDate(0, 0, 6).Weekday()},
			},
		},
		{
			Name: "simple",
			Input: []ReviewsAtTime{
				{Count: 2, Time: todayStart.Add(time.Hour * 18)},
				{Count: 1, Time: todayStart.Add(time.Hour * 19)},
				{Count: 3, Time: todayStart.Add(time.Hour * 22)},
				{Count: 4, Time: todayStart.AddDate(0, 0, 2).Add(time.Hour * 15)},
				{Count: 4, Time: todayStart.AddDate(0, 0, 3).Add(time.Hour * 12)},
			},
			Result: []DayForecast{
				{
					Weekday:    todayStart.Weekday(),
					TotalCount: 6,
					Hourly: []HourForecast{
						{
							Time:            todayStart.Add(time.Hour * 18),
							Count:           2,
							CumulativeCount: 2,
						},
						{
							Time:            todayStart.Add(time.Hour * 19),
							Count:           1,
							CumulativeCount: 3,
						},
						{
							Time:            todayStart.Add(time.Hour * 22),
							Count:           3,
							CumulativeCount: 6,
						},
					},
				},
				{Weekday: todayStart.AddDate(0, 0, 1).Weekday()},
				{
					Weekday:    todayStart.AddDate(0, 0, 2).Weekday(),
					TotalCount: 4,
					Hourly: []HourForecast{
						{
							Time:            todayStart.AddDate(0, 0, 2).Add(time.Hour * 15),
							Count:           4,
							CumulativeCount: 10,
						},
					},
				},
				{
					Weekday:    todayStart.AddDate(0, 0, 3).Weekday(),
					TotalCount: 4,
					Hourly: []HourForecast{
						{
							Time:            todayStart.AddDate(0, 0, 3).Add(time.Hour * 12),
							Count:           4,
							CumulativeCount: 14,
						},
					},
				},
				{Weekday: todayStart.AddDate(0, 0, 4).Weekday()},
				{Weekday: todayStart.AddDate(0, 0, 5).Weekday()},
				{Weekday: todayStart.AddDate(0, 0, 6).Weekday()},
			},
		},
		{
			Name: "next review isn't today",
			Input: []ReviewsAtTime{
				{Count: 5, Time: todayStart.AddDate(0, 0, 2).Add(time.Hour * 9)},
				{Count: 1, Time: todayStart.AddDate(0, 0, 2).Add(time.Hour * 12)},
				{Count: 9, Time: todayStart.AddDate(0, 0, 4).Add(time.Hour * 16)},
				{Count: 7, Time: todayStart.AddDate(0, 0, 6).Add(time.Hour * 18)},
			},
			Result: []DayForecast{
				{Weekday: now.Weekday()},
				{Weekday: now.AddDate(0, 0, 1).Weekday()},
				{
					Weekday:    now.AddDate(0, 0, 2).Weekday(),
					TotalCount: 6,
					Hourly: []HourForecast{
						{
							Time:            todayStart.AddDate(0, 0, 2).Add(time.Hour * 9),
							Count:           5,
							CumulativeCount: 5,
						},
						{
							Time:            todayStart.AddDate(0, 0, 2).Add(time.Hour * 12),
							Count:           1,
							CumulativeCount: 6,
						},
					},
				},
				{Weekday: now.AddDate(0, 0, 3).Weekday()},
				{
					Weekday:    now.AddDate(0, 0, 4).Weekday(),
					TotalCount: 9,
					Hourly: []HourForecast{
						{
							Time:            todayStart.AddDate(0, 0, 4).Add(time.Hour * 16),
							Count:           9,
							CumulativeCount: 15,
						},
					},
				},
				{Weekday: now.AddDate(0, 0, 5).Weekday()},
				{
					Weekday:    now.AddDate(0, 0, 6).Weekday(),
					TotalCount: 7,
					Hourly: []HourForecast{
						{
							Time:            todayStart.AddDate(0, 0, 6).Add(time.Hour * 18),
							Count:           7,
							CumulativeCount: 22,
						},
					},
				},
			},
		},
		{
			Name: "with reviews exceeding week",
			Input: []ReviewsAtTime{
				{Count: 5, Time: todayStart.AddDate(0, 0, 2).Add(time.Hour * 9)},
				{Count: 7, Time: todayStart.AddDate(0, 0, 7).Add(time.Hour * 18)},
			},
			Result: []DayForecast{
				{Weekday: now.Weekday()},
				{Weekday: now.AddDate(0, 0, 1).Weekday()},
				{Weekday: now.AddDate(0, 0, 2).Weekday(),
					TotalCount: 5,
					Hourly: []HourForecast{
						{
							Time:            todayStart.AddDate(0, 0, 2).Add(time.Hour * 9),
							Count:           5,
							CumulativeCount: 5,
						},
					},
				},
				{Weekday: now.AddDate(0, 0, 3).Weekday()},
				{Weekday: now.AddDate(0, 0, 4).Weekday()},
				{Weekday: now.AddDate(0, 0, 5).Weekday()},
				{Weekday: now.AddDate(0, 0, 6).Weekday()},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			result := NewForecast(tc.Input, loc)

			assert.Equal(t, tc.Result, result)
		})
	}
}
