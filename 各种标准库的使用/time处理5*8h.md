
>
>
>输入5\*8h 7\*24h、开始时间、结束时间、返回时间间隔
```python
package main

import (
   "fmt"
   "strconv"
   "strings"
   "time"
)

// 是否在工作时间
func isWorkTime(i int) bool {
   var (
      WorkTimeSlice = []int{9, 10, 11, 12, 13, 14, 15, 16, 17, 18}
   )
   for _, v := range WorkTimeSlice {
      if v == i {
         return true
      }
   }
   return false
}

// 是否为周末
func isWeekDayTime(i time.Weekday) bool {
   var (
      weekDaySlice = []time.Weekday{
         time.Sunday,
         time.Saturday,
      }
   )
   for _, v := range weekDaySlice {
      if v == i {
         return true
      }
   }
   return false
}

// 解析字符串
func pasTime(workTime string) (day int, hour int, err error) {
   workTime = strings.ReplaceAll(workTime, "h", "")
   workTimeSlice := strings.Split(workTime, "*")
   if day, err = strconv.Atoi(workTimeSlice[0]); err != nil {
      return
   }
   if hour, err = strconv.Atoi(workTimeSlice[1]); err != nil {
      return
   }
   return
}

// 返回日期工作时间点
func workTimeObj(timeObj time.Time, workTime int) (workTimeObj time.Time) {
   loc, err := time.LoadLocation("Asia/Shanghai")
   if err != nil {
      fmt.Println(err)
      return
   }
   workTimeFormat := fmt.Sprintf("%v-%v-%v %v:00:00", timeObj.Year(), int(timeObj.Month()), timeObj.Day(), workTime)
   workTimeObj, _ = time.ParseInLocation("2006-1-2 15:04:05", workTimeFormat, loc)
   return
}

func analysisTime(day, hour int, startTime, endTime time.Time) (intervalTime time.Duration) {
   // 1. 7* 24小时制度
   if hour == 24 && day == 7 {
      return endTime.Sub(startTime)
   }

   if startTime.Hour() < 9 {
      startTime = workTimeObj(startTime, 9)
   }
   if startTime.Hour() > 18 {
      startTime = workTimeObj(startTime, 9).Add(24 * time.Hour)
   }
   if endTime.Hour() < 9 {
      endTime = workTimeObj(endTime, 18).Add(-24 * time.Hour)
   }
   if endTime.Hour() > 18 {
      endTime = workTimeObj(endTime, 18)
   }

   // 3. 创建工单时间 和 响应工单时间 不同天,start end都为工作时间
   for {
      if startTime.Month() == endTime.Month() {
         if startTime.Day() > endTime.Day() {
            break
         }
      }
      //2. 创建工单时间 和 响应工单时间 同一天
      if startTime.Day() != endTime.Day() {
			if !isWeekDayTime(startTime.Weekday()) { // 不是周末
				if endTime.Day()-startTime.Day() > 1 {
					intervalTime += 9 * time.Hour
				}
			}
         startTime = startTime.Add(24 * time.Hour)
         continue
      }
      if endTime.After(startTime) {
         if isWorkTime(startTime.Hour()) && isWorkTime(endTime.Hour()) { //工作时间
            intervalTime += endTime.Sub(startTime)
         }
      } else {
         workStartTimeObj := workTimeObj(endTime, 9)
         intervalTime += endTime.Sub(workStartTimeObj)
         workEndTimeObj := workTimeObj(endTime, 18)
         intervalTime += workEndTimeObj.Sub(startTime) // 间隔时间
      }
      break
   }
   return

}

func main() {
   loc, _ := time.LoadLocation("Asia/Shanghai")
   startTime, _ := time.ParseInLocation("2006-1-2 15:04:05", "2020-3-26 17:45:00", loc) // 创建时间
   endTime, _ := time.ParseInLocation("2006-1-2 15:04:05", "2020-3-27 8:15:00", loc)    //结束时间
   day, hour, _ := pasTime("5*8h")
   abs := analysisTime(day, hour, startTime, endTime)
   fmt.Println(abs)

}
```