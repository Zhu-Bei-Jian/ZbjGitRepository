package gameutil

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

const (
	goBornDate       = "2006-01-02"          //长日期格式
	goBornShortDate  = "06-01-02"            //短日期格式
	goBornTimes      = "15:04:05"            //长时间格式
	goBornShortTime  = "15:04"               //短时间格式
	goBornDateTime   = "2006-01-02 15:04:05" //日期时间格式
	goBornDateString = "20060102"
	//goBornDateStringTime = "20060102150405"
)

var monthIntMap = map[string]int{
	"January":   1,
	"February":  2,
	"March":     3,
	"April":     4,
	"May":       5,
	"June":      6,
	"July":      7,
	"August":    8,
	"September": 9,
	"October":   10,
	"November":  11,
	"December":  12,
}

// GetCurrentTime 获取时间.
func GetCurrentTime() time.Time {
	return time.Now()
}

// GetCurrentTimestamp 获取时间戳.
func GetCurrentTimestamp() int64 {
	return GetCurrentTime().Unix()
}

// GetCurrentMicroTimestamp 获取毫秒时间戳.
func GetCurrentMicroTimestamp() int64 {
	return GetCurrentTime().UnixNano() / 1e6
}

// GetCurrentNanoTimestamp 获取纳秒时间戳.
func GetCurrentNanoTimestamp() int64 {
	return GetCurrentTime().UnixNano()
}

// GetTodayZeroTime 获取本地时间的零点时间结构体.
func GetTodayZeroTime() time.Time {
	timeStr := GetCurrentTime().Format(goBornDate)
	zeroTime, _ := time.ParseInLocation(goBornDate, timeStr, time.Local)
	return zeroTime
}

func GetNextDayZeroTime() time.Time {
	timeStr := GetCurrentTime().Add(time.Hour * 24).Format(goBornDate)
	t, _ := time.ParseInLocation(goBornDate, timeStr, time.Local)
	return t
}

// GetNextWeekZeroTime... 获取下一周零点时间戳
func GetNextWeekZeroTime() time.Time {
	timeStr := GetTargetDayWeekStartZeroTime(GetCurrentTime()).Add(time.Hour * 7 * 24).Format(goBornDate)
	t, _ := time.ParseInLocation(goBornDate, timeStr, time.Local)
	return t
}

func GetWeekStartZeroTime() int64 {
	return GetTargetDayWeekStartZeroTime(GetCurrentTime()).Unix()
}

func GetMonthStartZeroTime() time.Time {
	return GetTargetDayMonthStartZeroTime(GetCurrentTime())
}

func GetNextMonthStartZeroTime() time.Time {
	return GetTargetDayNextMonthStartZeroTime(GetCurrentTime())
}

func GetTargetDayNextMonthStartZeroTime(t time.Time) time.Time {
	year, month, _ := t.Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	NextMonth := thisMonth.AddDate(0, 1, 0)
	return NextMonth
}

func GetTargetDayMonthStartZeroTime(t time.Time) time.Time {
	year, month, _ := t.Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	return thisMonth
}

// GetTargetDayWeekStartZeroTime 获取指定时间的周开始零点时间.
func GetTargetDayWeekStartZeroTime(t time.Time) time.Time {
	weekDay := GetWeekDay(t)
	if weekDay == 0 {
		weekDay = 7
	}
	year, month, day := t.Date()
	thisWeek := time.Date(year, month, day-weekDay+1, 0, 0, 0, 0, time.Local)
	return thisWeek
}

// ParseTimeStr2Timestamp 把相对时间 xx:xx:xx 转换成当天对应的时间戳.
func ParseTimeStr2Timestamp(timeStr string) (int64, error) {
	str := fmt.Sprintf("%s %s", GetCurrentTime().Format(goBornDate), timeStr)
	newTime, err := time.ParseInLocation(goBornDateTime, str, time.Local)
	if err != nil {
		return 0, err
	}
	return newTime.Unix(), nil
}

// ParseTimeStr2Second 把相对时间 xx:xx:xx 转换成相应的秒数.
func ParseTimeStr2Second(timeStr string) (int64, error) {
	s, err := ParseTimeStr2Timestamp(timeStr)
	if err != nil {
		return 0, err
	}
	s -= GetTodayZeroTime().Unix()
	return s, nil
}

// ParseDate2Timestamp 把相对时间 2018-01-01 转换成 时间戳.
func ParseDate2Timestamp(timeStr string) (int64, error) {
	// 转化成时间格式
	loc, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation(goBornDate, timeStr, loc)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// ParseDatetime2Timestamp 把相对时间 2018-01-01 00:00:00 转换成 时间戳.
func ParseDatetime2Timestamp(timeStr string) (int64, error) {
	// 转化成时间格式
	loc, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation(goBornDateTime, timeStr, loc)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func ParseDateTime2Time(timeStr string) (time.Time, error) {
	// 转化成时间格式
	loc, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation(goBornDateTime, timeStr, loc)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

//yyyyMMddHHmmss
func ParseFormatDateTime2Time(timeStr string, format string) (time.Time, error) {
	// 转化成时间格式
	loc, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation(format, timeStr, loc)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// ParseTimestamp2Time 把相对时间戳转换成相应的时间.
func ParseTimestamp2Time(sec int64) time.Time {
	return time.Unix(sec, 0)
}

// ParseTimestamp2String 把相对时间戳转换成相应的字符串时间.
func ParseTimestamp2String(sec int64) string {
	return time.Unix(sec, 0).Format("2006-01-02 15:04:05")
}

// ParseTime2Timestamp 把相对时间转换成相应的时间戳.
func ParseTime2Timestamp(t time.Time) int64 {
	return t.UTC().Unix()
}

func ParseTime2String(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func GetCurrentDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetCurrentDateTimeHHSS() string {
	return time.Now().Format("20060102150405")
}

// IsInTimeNow 判断当前时间是否在给定时间范围内.
func IsInTimeNow(start, end string) bool {
	currentTime := GetCurrentTimestamp()
	// 判断是否有指定时间
	if start == "" || end == "" {
		return false
	}
	// 解析时间字符串
	startTime, err := ParseDatetime2Timestamp(start)
	if err != nil {
		return false
	}
	endTime, err := ParseDatetime2Timestamp(end)
	if err != nil {
		return false
	}
	// 判断是否在时间范围内
	if currentTime >= startTime && currentTime <= endTime {
		return true
	}
	return false
}

// GetMonth 获取指定时间所属月份.
func GetMonth(t time.Time) int {
	return int(t.Month())
}

// GetWeekDay 获取指定时间所属星期.
func GetWeekDay(t time.Time) int {
	return int(t.Weekday())
}

// GetWeekCountToday 获取从1970-01-01 00:00:00到现在过了多少周.!!!以周1为第一天
func GetWeekCountToday() int {
	t := GetCurrentTimestamp() + 8*3600
	t += 86400 * 3 // 1970-01-01 周四
	weekCount := math.Floor(float64(t / (86400 * 7)))
	return int(weekCount)
}

// GetTargetDayWeekCount 获取从1970-01-01 00:00:00到指定时间过了多少周.
func GetTargetDayWeekCount(second int64) int {
	t := GetTargetDayZeroTime(IntToTime(second)).Unix() + 8*3600
	t += 86400 * 3 // 1970-01-01 周四
	weekCount := math.Floor(float64(t / (86400 * 7)))
	return int(weekCount)
}

// GetTargetDayLastMonthStartZeroTime 获取指定时间的上个月开始零点时间.
func GetTargetDayLastMonthStartZeroTime(t time.Time) time.Time {
	year, month, _ := t.Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastMonth := thisMonth.AddDate(0, -1, 0)
	return lastMonth
}

// GetTargetDayLastMonthEndZeroTime 获取指定时间的上个月结束零点时间.
func GetTargetDayLastMonthEndZeroTime(t time.Time) time.Time {
	year, month, _ := t.Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastMonth := thisMonth.AddDate(0, 0, -1)
	return lastMonth
}

// GetLastMonthStartZeroTime 获取当前时间的上个月开始零点时间.
func GetLastMonthStartZeroTime() time.Time {
	return GetTargetDayLastMonthStartZeroTime(GetCurrentTime())
}

// GetLastMonthEndZeroTime 获取当前时间的上个月结束零点时间.
func GetLastMonthEndZeroTime(t time.Time) time.Time {
	return GetTargetDayLastMonthEndZeroTime(GetCurrentTime())
}

// GetTargetDayLastWeekStartZeroTime 获取指定时间的上周开始零点时间.
func GetTargetDayLastWeekStartZeroTime(t time.Time) time.Time {
	weekDay := GetWeekDay(t)
	year, month, day := t.Date()
	thisWeek := time.Date(year, month, day-weekDay, 0, 0, 0, 0, time.Local)
	lastWeek := thisWeek.AddDate(0, 0, -7)
	return lastWeek
}

// GetTargetDayLastWeekEndZeroTime 获取指定时间的上周结束零点时间.
func GetTargetDayLastWeekEndZeroTime(t time.Time) time.Time {
	weekDay := GetWeekDay(t)
	year, month, day := t.Date()
	thisWeek := time.Date(year, month, day-weekDay, 0, 0, 0, 0, time.Local)
	lastWeek := thisWeek.AddDate(0, 0, -1)
	return lastWeek
}

// GetLastWeekStartZeroTime 获取当前时间的上周开始零点时间.
//func GetLastWeekStartZeroTime() time.Time {
//	return GetTargetDayLastWeekStartZeroTime(GetCurrentTime())
//}
//
//// GetLastWeekEndZeroTime 获取当前时间的上周结束零点时间.
//func GetLastWeekEndZeroTime() time.Time {
//	return GetTargetDayLastWeekEndZeroTime(GetCurrentTime())
//}

func GetLastWeekStartTime() time.Time {
	timeStr := GetTargetDayWeekStartZeroTime(GetCurrentTime()).Add(-time.Hour * 7 * 24).Format(goBornDate)
	t, _ := time.ParseInLocation(goBornDate, timeStr, time.Local)
	return t
}

func GetLastWeekEndTime() time.Time {
	timeStr := GetTargetDayWeekStartZeroTime(GetCurrentTime()).Format(goBornDate)
	t, _ := time.ParseInLocation(goBornDate, timeStr, time.Local)
	return t
}

func GetCurrentWeekDay() int32 {
	return int32(time.Now().Weekday())
}

// GetTargetDayZeroTime 获取指定时间的零点.
func GetTargetDayZeroTime(t time.Time) time.Time {
	timeStr := t.Format(goBornDate)
	zeroTime, _ := time.ParseInLocation(goBornDate, timeStr, time.Local)
	return zeroTime
}

// GetTargetDayZeroTimeInt 获取指定时间零点时间戳.
func GetTargetDayZeroTimeInt(t time.Time) int64 {
	return GetTargetDayZeroTime(t).Unix()
}

// GetTargetDaySecond2ZeroTimeInt 获取指定时间零点时间戳.
func GetTargetDaySecond2ZeroTimeInt(seconds int64) int64 {
	t := IntToTime(seconds)
	return GetTargetDayZeroTime(t).Unix()
}

// IntToTime 将时间戳转换为Time结构.
func IntToTime(seconds int64) time.Time {
	return time.Unix(seconds, 0)
}

// DateStringToInt 将长日期格式的字符串转换为本地时间戳.
func DateStringToInt(date string) int64 {
	now, _ := time.ParseInLocation(goBornDate, date, time.Local)
	return now.Unix()
}

// GetTodayDateString 获取本地时间 年月日 格式的字符串.
func GetTodayDateString() string {
	return GetCurrentTime().Format(goBornDateString)
}

// GetTargetDayDateString 获取指定时间 年月日 格式的字符串.
func GetTargetDayDateString(t time.Time) string {
	return t.Format(goBornDateString)
}

// GetCurrentYear 获取当前年份
func GetCurrentYear() int {
	return time.Now().Year()
}

// GetCurrentYear 获取当前月份
func GetCurrentMonth() int {
	return monthIntMap[time.Now().Month().String()]
}

// GetCurrentDate 获取当天年月日
func GetCurrentDate() string {
	return time.Now().Format("20060102")
}

//// GetCountDay 获取从1970-01-01 00:00:00到现在过了多少天.cst时间多8个小时
func GetCountDayNow() int32 {
	return GetCountDay(GetCurrentTimestamp())
}

func GetCountDay(seconds int64) int32 {
	return GetDayIndexAfter2000ByUnix(seconds)
	//count := (seconds + 8*3600) / (60 * 60 * 24)
	//return int32(count)
}

func GetCurrentDayIndexAfter2000() int32 { //180101 18年1月1日
	sDayNow := time.Now().Format("20060102")
	dayNow, _ := strconv.Atoi(sDayNow)
	if dayNow > 20000000 {
		return int32(dayNow - 20000000)
	}
	return 0
}

//20060102
func GetTargetDateID(t time.Time) int32 {
	sDayNow := t.Format("20060102")
	dayNow, _ := strconv.Atoi(sDayNow)
	return int32(dayNow)
}

//0102
func GetTargetMonthDayID(t time.Time) int32 {
	sDayNow := t.Format("0102")
	dayNow, _ := strconv.Atoi(sDayNow)
	return int32(dayNow)
}

func GetTargetYearMonthID(t time.Time) int32 {
	sDayNow := t.Format("200601")
	dayNow, _ := strconv.Atoi(sDayNow)
	return int32(dayNow)
}

//01
func GetTargetMonth(t time.Time) int32 {
	sDayNow := t.Format("01")
	dayNow, _ := strconv.Atoi(sDayNow)
	return int32(dayNow)
}

func GetDayIndexAfter2000ByUnix(sec int64) int32 { //180101 18年1月1日
	if sec == 0 {
		return 0
	}
	t := time.Unix(sec, 0)
	sDayNow := t.Format("20060102")
	dayNow, _ := strconv.Atoi(sDayNow)
	if dayNow > 20000000 {
		return int32(dayNow - 20000000)
	}
	return 0
}
func GetMonthIndexAfter2000ByUnix(sec int64) int32 { //1801 18年1月
	return GetDayIndexAfter2000ByUnix(sec) / 100
}

func GetCurrentWeekIndexAfter2000ByUnix(sec int64) int32 {
	return int32(GetTargetDayWeekCount(sec))
}

func GetNowFormatHour() int32 {
	sHourNow := time.Now().Format("2006010215")
	hourNow, _ := strconv.Atoi(sHourNow)
	return int32(hourNow)
}
func GetCurrentMonthIndexAfter2000() int32 { //1801 18年1月
	return GetCurrentDayIndexAfter2000() / 100
}
func GetCurrentWeekIndexAfter2000() int32 {
	return GetCurrentWeekIndexAfter2000ByUnix(GetCurrentTimestamp())
}

/*获得以特定开始时间的 几天为一个周期的 天索引：这个周期的第一天时间
*startDayTick 这个值要是0点0分的秒值 0以1970101开始
 */
func GetNDayIndexAfter2000(startDayTick int64, tick int64, dayN int64) int32 {
	if startDayTick == 0 {
		tBase, e := time.ParseInLocation(goBornDateTime, "1970-01-01 00:00:00", time.Local)
		if e == nil {
			startDayTick = tBase.Unix()
		}
	}
	if tick <= startDayTick {
		return 0
	}
	diff := tick - startDayTick
	twoDayBeginTick := startDayTick + diff - diff%(86400*dayN)
	return GetDayIndexAfter2000ByUnix(twoDayBeginTick)
}

// GetCurrentDay 获取当前日
func GetCurrentDay() uint32 {
	return uint32(time.Now().Day())
}

// GetMonthDay 获取指定时间的天数
func GetMonthDay(year int, month int) int {
	days := 0
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			days = 30

		} else {
			days = 31
		}
	} else {
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			days = 29
		} else {
			days = 28
		}
	}
	return days
}

func GetCurrentMonthDay() int {
	return GetMonthDay(GetCurrentYear(), GetCurrentMonth())
}

// GetTodayNHourTimeInt 获取本地之间今天第N小时的时间戳
func GetTodayNHourTimeInt(n int) int64 {
	return GetTodayZeroTimeInt() + int64(n)*int64(time.Hour/time.Second)
}

//获取指定时间当天第N小时的时间戳
func GetTargetNHourTimeInt(t time.Time, n int) int64 {
	return GetTargetDayZeroTimeInt(t) + int64(n)*int64(time.Hour/time.Second)
}

// GetTodayZeroTimeInt 获取本地时间的零点时间戳，单位秒
func GetTodayZeroTimeInt() int64 {
	return GetTodayZeroTime().Unix()
}

// DiffDay 获取两个时间戳的相隔天数,以hour为分界点
func DiffDayNHour(left, right int64, hour int) int {
	leftTime := IntToTime(left)
	rightTime := IntToTime(right)

	var diff int64
	if left > right {
		diff = GetTargetLastestHourTimeInt(leftTime, hour) - GetTargetLastestHourTimeInt(rightTime, hour)
	} else {
		diff = GetTargetLastestHourTimeInt(rightTime, hour) - GetTargetLastestHourTimeInt(leftTime, hour)
	}

	return int(diff / (24 * 3600))
}

//获取指定时间之前的离hour最近的时间戳
func GetTargetLastestHourTimeInt(t time.Time, hour int) int64 {
	//以五点为例，在五点时间之后的，以当天五点为时间戳，在五点时间之前的，以昨天五点为时间戳
	if t.Hour() >= hour {
		return GetTargetNHourTimeInt(t, hour)
	} else {
		return GetTargetNHourTimeInt(t, hour) - 24*3600
	}
}

//打印时间流逝
func PrintElaspe(t time.Time) {
	fmt.Println(time.Since(t))
}

func PrintElaspeLog(title string, t time.Time) {
	fmt.Println(title, ":", time.Since(t))
}

//当前是19:59返回 1959
func GetCurrHM() int32 {
	now := time.Now()
	hm := now.Hour()*100 + now.Minute()
	return int32(hm)
}

func GetCountDay1970(t int64) int32 {
	t += 8 * 3600
	count := math.Floor(float64(t / (60 * 60 * 24)))
	return int32(count)
}
