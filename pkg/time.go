package pkg

import (
	"fmt"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/common"
	"math"
	"strings"
	"time"
)

// ParseTime 解析字幕时间字符串，这里可能小数点后面有 2-4 位
func ParseTime(inTime string) (time.Time, error) {

	parseTime, err := time.Parse(common.TimeFormatPoint2, inTime)
	if err != nil {
		parseTime, err = time.Parse(common.TimeFormatPoint3, inTime)
		if err != nil {
			parseTime, err = time.Parse(common.TimeFormatPoint4, inTime)
		}
	}
	return parseTime, err
}

func TimeNumber2Time(inputTimeNumber float64) time.Time {
	newTime := time.Time{}.Add(time.Duration(inputTimeNumber * math.Pow10(9)))
	return newTime
}

func Time2SecondNumber(inTime time.Time) float64 {
	outSecond := 0.0
	outSecond += float64(inTime.Hour() * 60 * 60)
	outSecond += float64(inTime.Minute() * 60)
	outSecond += float64(inTime.Second())
	outSecond += float64(inTime.Nanosecond()) / 1000 / 1000 / 1000

	return outSecond
}

func Time2Duration(inTime time.Time) time.Duration {
	return time.Duration(Time2SecondNumber(inTime) * math.Pow10(9))
}

// Time2SubTimeString 时间转字幕格式的时间字符串
func Time2SubTimeString(inTime time.Time, timeFormat string) string {
	/*
		这里进行时间转字符串的时候有一点比较特殊
		正常来说输出的格式是类似 15:04:05.00
		那么有个问题，字幕的时间格式是 0:00:12.00， 小时，是个数，除非有跨度到 20 小时的视频，不然小时就应该是个数
		这就需要一个额外的函数去处理这些情况
	*/
	outTimeString := inTime.Format(timeFormat)
	if inTime.Hour() > 9 {
		// 小时，两位数
		return outTimeString
	} else {
		// 小时，一位数
		items := strings.SplitN(outTimeString, ":", -1)
		if len(items) == 3 {

			outTimeString = strings.Replace(outTimeString, items[0], fmt.Sprintf("%d", inTime.Hour()), 1)
			return outTimeString
		}

		return outTimeString
	}
}

func UnixTime2Time(unixTime int64) time.Time {
	return time.Unix(unixTime, 0)
}

// SecondsToHMS 将秒转换为 00:01:54,830 格式
func SecondsToHMS(seconds float64) string {
	hours := int(math.Floor(seconds / 3600))
	minutes := int(math.Floor((seconds - float64(hours)*3600) / 60))
	secs := int(math.Floor(seconds - float64(hours)*3600 - float64(minutes)*60))
	msecs := int(math.Floor((seconds - float64(hours)*3600 - float64(minutes)*60 - float64(secs)) * 1000))
	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, secs, msecs)
}
