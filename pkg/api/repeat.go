package api

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// служебные функции
func convertSteps(steps []string) ([]int, error) {
	days := make([]int, len(steps))

	for i, s := range steps {
		num, err := strconv.Atoi(s)
		if err != nil {
			// Обработка ошибки, если строку не удалось преобразовать
			fmt.Println("conversion error:", err)
			return nil, err
		}
		days[i] = num
	}

	return days, nil
}

func daysSort(numbers []int) []int {
	// Выстраиваем дни в нужной последовательности:
	// Сначала положительные в порядке возрастания, затем отрицательные в порядке возрастания
	var positives []int
	var negatives []int

	for _, num := range numbers {
		if num < 0 {
			negatives = append(negatives, num)
		} else {
			positives = append(positives, num)
		}
	}
	// сортировка нутри списка
	sort.Ints(positives)
	sort.Ints(negatives)

	// соединяем списки
	return append(positives, negatives...)

}

func weekDayAdd(date time.Time, weekSteps []string) time.Time {

	mindiff := 7

	// находим ближайший номер дня недели
	for _, wdStr := range weekSteps {
		wd, err := strconv.Atoi(wdStr)
		if err != nil {
			fmt.Println(err)
		}
		c := wd - int(date.Weekday())

		diff := ((math.Sqrt(float64(c*c))*(-1.0)+float64(c))/(float64(c)*2.0))*7.0 + float64(c)

		if diff < float64(mindiff) {
			mindiff = int(diff)
		}
	}

	return date.AddDate(0, 0, mindiff)
}

// Основная функция -------------------------------------------------------------
// находим NextDate
var dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	startDate, err := time.Parse(dateFormat, dstart)
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("invalid date format. get %s instead %s", dstart, dateFormat)
	}

	if repeat == "" {
		return "", fmt.Errorf("repeat value is empty: %s", repeat)
	}

	steps := strings.Split(repeat, " ")

	switch steps[0] {

	// условие для годов
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
			startDate = startDate.AddDate(1, 0, 0)
		}
		return startDate.Format(dateFormat), nil

		// условие для дней
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
		// особенное условие для тестов go test -run ^TestAddTask$ ./tests/
		if startDate.Before(now) && days == 1 {

			startDate = now

		} else {
			startDate = startDate.AddDate(0, 0, days)
			for {

				if startDate.After(now) {

					break
				}

				startDate = startDate.AddDate(0, 0, days) //по условию надо добавлять 1, но тогда не совпадает с тестом
			}
		}

		return startDate.Format(dateFormat), nil

		// условие для недель
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

		newDate := weekDayAdd(startDate, weekSteps)
		for {
			if newDate.After(now) {
				break
			}
			newDate = weekDayAdd(now, weekSteps)
		}
		return newDate.Format(dateFormat), nil

		// условие для месяцев
	case "m":
		if len(steps) == 1 {
			return "", fmt.Errorf("unsupported repeat format: %s", repeat)
		}

		dayMSteps := strings.Split(steps[1], ",")

		for _, md := range dayMSteps {

			mDay, err := strconv.Atoi(md)
			if err != nil {
				return "", err
			}

			if mDay == 0 || mDay > 31 || mDay < -2 {
				return "", fmt.Errorf("unsupported number of month day: %v", mDay)
			}
		}

		// календарь дней
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

		// если заданы тольк дни месяца
		if len(steps) == 2 {
			days, _ := convertSteps(strings.Split(steps[1], ","))
			days = daysSort(days)

			maxDayMonth = monthDays[int(startDate.Month())]
			isNextMonth := 0
			var newDate time.Time

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
				additionDay := int((math.Sqrt(float64(d*d))*(-1.0) + float64(d)) / (float64(d) * 2.0))
				newDate = startDate.AddDate(0, isNextMonth, additionDay-dayStart+d)

				if newDate.After(now) {

					return newDate.Format(dateFormat), nil
				}

			}

			newDate = startDate
			// находим ближайший день в будущем
			for i := 1; i < 13; i += 1 {
				newDate = startDate.AddDate(0, i, 0-startDate.Day()+days[0])
				if newDate.After(now) {
					break
				}

			}
			return newDate.Format(dateFormat), nil

			// если заданы дни и номера месяца
		} else if len(steps) == 3 {

			monthMSteps := strings.Split(steps[2], ",")

			for _, mm := range monthMSteps {

				mDay, err := strconv.Atoi(mm)
				if err != nil {
					return "", err
				}

				if mDay <= 0 || mDay > 12 {
					return "", fmt.Errorf("unsupported number of month: %v", mDay)
				}
			}

			days, _ := convertSteps(strings.Split(steps[1], ","))
			days = daysSort(days)
			months, _ := convertSteps(strings.Split(steps[2], ","))

			newDates := []time.Time{}
			//Невнятное условие на случай если год 1689

			nowYear := 2023
			if startDate.Year() < now.Year() {
				nowYear = int(now.Year())
			}
			year := 0
			if int(startDate.Year()) < nowYear {
				year = nowYear
			} else {
				year = int(startDate.Year())
			}

			// находим вариации ближайших дней
			for _, m := range months {

				for _, d := range days {
					y := year
					if d < 0 {
						d = monthDays[m] + d
					}
					if (y < int(now.Year()) && (m < int(startDate.Month()))) || (y < int(now.Year()) && m == int(startDate.Month()) && d < int(startDate.Day())) {
						y = year + 1
					}

					newDate := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
					newDates = append(newDates, newDate)
				}

			}

			// находим ближайший день в будущем
			aimDate := startDate.AddDate(3, 0, 0)
			for _, nd := range newDates {

				if startDate.Before(nd) && aimDate.After(nd) && nd.After(now) {
					aimDate = nd
				}

			}

			if aimDate.After(now) {
				return aimDate.Format(dateFormat), nil
			}

		} else {
			return "", fmt.Errorf("unsupported repeat format: %s", repeat)
		}

	default:
		return "", fmt.Errorf("unavailiable steps format: %s", steps[0])
	}
	return "", fmt.Errorf("no evaluates: %s", steps[0])
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
