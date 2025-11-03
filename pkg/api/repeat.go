package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// NextDate
func Sort(steps []string) ([]int, error) {
	days := make([]int, len(steps))

	for i, s := range steps {
		num, err := strconv.Atoi(s)
		if err != nil {
			// Обработка ошибки, если строку не удалось преобразовать
			fmt.Println("Ошибка преобразования:", err)
			return nil, err
		}
		days[i] = num
	}

	return days, nil
}

func daysSort(numbers []int) []int {
	var positives []int
	var negatives []int

	for _, num := range numbers {
		if num < 0 {
			negatives = append(negatives, num)
		} else {
			positives = append(positives, num)
		}
	}

	// Шаг 2: Сортируем списки
	sort.Ints(positives)                              // Положительные в порядке возрастания
	sort.Sort(sort.Reverse(sort.IntSlice(negatives))) // Отрицательные в порядке убывания

	// Шаг 3: Объединяем два списка
	return append(positives, negatives...)

}

var dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	startDate, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format. get %s instead %s", dstart, dateFormat)
	}

	if repeat == "" {
		return "", fmt.Errorf("repeat value is empty: %s", repeat)
	}

	steps := strings.Split(repeat, " ")

	fmt.Println(steps[0])

	switch steps[0] {
	case "y":
		if len(steps) == 1 {
			//
			startDate = startDate.AddDate(1, 0, 0)
		} else {
			yearStep, err := strconv.Atoi(steps[1])
			if err != nil {
				return "", fmt.Errorf("unsupported repeat format: %s", repeat)
			}
			startDate = startDate.AddDate(yearStep, 0, 0)
		}
		for {

			if startDate.After(now) {
				break
			}
			startDate = now.AddDate(0, 0, 1)
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

		if days <= 0 || days > 400 {
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

		weekSteps := strings.Split(steps[1], ",")

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
		febDay := 28
		maxDayMonth := 0
		if startDate.Year()%4 == 0 {
			febDay = 29
		}
		monthDays := map[int]int{
			1:  31,
			2:  febDay,
			3:  31,
			4:  30,
			5:  31,
			6:  30,
			7:  31,
			8:  31,
			9:  30,
			10: 31,
			11: 30,
			12: 31,
		}
		dayStart := startDate.Day()

		if len(steps) == 2 {
			days, _ := Sort(strings.Split(steps[1], ","))
			days = daysSort(days)

			maxDayMonth = monthDays[int(startDate.Month())]
			isNextMonth := 0

			for _, d := range days {
				isNextMonth = 0

				if startDate.Day() >= d && d <= maxDayMonth {
					if d > monthDays[int(startDate.Month())+1] {
						isNextMonth = 2
					} else {
						isNextMonth = 1
					}

				} else if d > maxDayMonth {
					isNextMonth = 1
				}

				newDate := startDate.AddDate(0, isNextMonth, 0-dayStart+d)
				if newDate.After(now) {

					return newDate.Format(dateFormat), nil
				}

			}

			for {
				if startDate.After(now) {
					break
				}
				startDate = now.AddDate(0, 0, 1)
			}

			return startDate.Format(dateFormat), nil

		} else if len(steps) == 3 {
			days, _ := Sort(strings.Split(steps[1], ","))
			days = daysSort(days)
			months, _ := Sort(strings.Split(steps[2], ","))

			newDates := []time.Time{}

			for _, m := range months {

				for _, d := range days {
					y := int(startDate.Year())
					if d < 0 {
						d = monthDays[m] + d
					}
					if (m < int(startDate.Month())) || (m == int(startDate.Month()) && d < int(startDate.Day())) {
						y = int(startDate.Year()) + 1
					}

					newDate := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
					newDates = append(newDates, newDate)
				}

			}
			aimDate := startDate.AddDate(3, 0, 0)
			for _, nd := range newDates {
				if startDate.Before(nd) && aimDate.After(nd) {
					aimDate = nd
				}

			}

			for {
				if aimDate.After(now) {
					break
				}
				aimDate = now.AddDate(0, 0, 1)
			}

			return aimDate.Format(dateFormat), nil

		} else {
			return "", fmt.Errorf("unsupported repeat format: %s", repeat)
		}

	default:
		return "", fmt.Errorf("unavailiable steps format: %s", steps[0])
	}

}

func writeJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to write JSON", http.StatusInternalServerError)
	}
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJson(w, map[string]string{"error": "method not allowed"})
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	if dateStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "date missing"})
		return
	}

	repeatStr := r.FormValue("repeat")
	if repeatStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "repeat missing"})
	}

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJson(w, map[string]string{"error": ""})
			return
		}
	}

	next, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte(next)); err != nil {
		log.Printf("failed to write response in /api/nextdate: %v", err)
	}
}
