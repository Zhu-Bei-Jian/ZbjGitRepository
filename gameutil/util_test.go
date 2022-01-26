package gameutil

import (
	"fmt"
	"testing"
)

func Test_Date(t *testing.T) {
	time1 := "2018-06-10 04:00:00"

	timestamp, _ := ParseDatetime2Timestamp(time1)
	//targetTime := ParseTimestamp2Time(timestamp)

	diffDay := DiffDayNHour(timestamp, GetCurrentTimestamp(), 5)

	if diffDay == 2 {
		t.Log(diffDay)
	} else {
		t.Error(diffDay)
	}

}

func Test_RandDiff(t *testing.T) {
	t.Log(GenDiffRandomNum(2, 5))
	t.Log(GenDiffRandomNum(8, 5))
}

func Test_Version(t *testing.T) {
	fmt.Println(CompareStringVersion("2.0.6.1", "2.0.6"))
	fmt.Println(IsVersionSmaller("2.0.6.1", "2.0.6"))
	fmt.Println(IsVersionGreater("2.0.6.1", "2.0.6"))
	fmt.Println(IsVersionSmaller("2.0.6.1", "2.0.6"))
	fmt.Println(IsVersionSmaller("2.0.6.1", "2.0.7"))
}
