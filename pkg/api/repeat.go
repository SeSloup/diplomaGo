package api

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// NextDate

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	dateFormat := "20060102"

	startDate, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format. get %s instead %s", dstart, dateFormat)
	}

	if repeat == "" {
		return "", fmt.Errorf("repeat value is empty: %s", repeat)
	}

	steps := strings.Split(repeat, " ")

	switch steps[0] {
	case "y":
		if len(steps) == 1 {
			return "", fmt.Errorf("unsupported repeat format: %s", repeat)
		}

		for {
			startDate = startDate.AddDate(1, 0, 0)
			if startDate.After(now) {
				break
			}
		}
		return startDate.Format(dateFormat), nil

	case "d":
		if len(steps) == 1 {
			return "", fmt.Errorf("unsupported repeat format: %s", repeat)
		}

		days, err := strconv.Atoi(steps[1])
		if err != nil {
			return "", err
		}

		if days <= 0 || days > 366 {
			return "", fmt.Errorf("unsupported count of days: %v", days)
		}

		startDate = startDate.AddDate(0, 0, days)
		for {

			if startDate.After(now) {
				break
			}
			startDate = now.AddDate(0, 0, 1)
		}
		return startDate.Format(dateFormat), nil

	case "w":
		if len(steps) == 1 {
			return "", fmt.Errorf("unsupported repeat format: %s", repeat)
		}

		weekSteps := strings.Split(steps[1], " ")

		for _, wd := range weekSteps {
			weekDay, err := strconv.Atoi(wd)
			if err != nil {
				return "", err
			}

			if weekDay <= 0 || weekDay > 7 {
				return "", fmt.Errorf("unsupported number of week day: %v", weekDay)
			}
		}

		dayOfWeek := startDate.Weekday() // Get the day of the week
		dayNumber := int(dayOfWeek)      // Convert to integer

		mindiff := 7
		//nearestWeekDay := 0

		// находим ближайший номер дня недели
		for _, wdStr := range weekSteps {
			wd, err := strconv.Atoi(wdStr)
			if err != nil {
				fmt.Println(err)
			}
			c := wd - dayNumber

			diff := ((math.Sqrt(float64(c*c))*(-1.0)+float64(c))/(float64(c)*2.0))*7.0 + float64(c)
			fmt.Println("diff", diff, "c", c)

			if diff < float64(mindiff) {
				mindiff = int(diff)
				//nearestWeekDay = wd
			}

		}
		startDate = startDate.AddDate(0, 0, mindiff)
		for {

			if startDate.After(now) {
				break
			}
			startDate = now.AddDate(0, 0, 1)
		}

		return startDate.Format(dateFormat), nil

		///

	case "m":

		///
		// ПРосчитать
		///
		if len(steps) == 2 {

		} else if len(steps) == 3 {

		} else {
			return "", fmt.Errorf("unsupported repeat format: %s", repeat)
		}

	default:
		return "", fmt.Errorf("unavailiable steps format: %s", steps[0])
	}
}
