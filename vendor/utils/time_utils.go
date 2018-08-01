package utils

import "time"

var LOCATIONZONE, _ = time.LoadLocation(LOCATION)

//create by cwj on 2017-10-17
// return now time by string
func Now() string {
	return time.Now().In(LOCATIONZONE).Format(TIME_FORMAT_1)
}

func ZeroTime() time.Time {
	a, _ := time.Parse(TIME_FORMAT_3, "0001-01-01")
	return a
}

func TodayWithout() string {
	return time.Now().In(LOCATIONZONE).Format("20060102")
}

func Today() string {
	return time.Now().In(LOCATIONZONE).Format("2006-01-02")
}

func GetToday() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, ZERO, ZERO, ZERO, ZERO, LOCATIONZONE)
}

func GetToday24() time.Time {
	return GetToday().Add(24 * time.Hour)
}
