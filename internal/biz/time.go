package biz

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	LocalDateTimeFormat string = "2006-01-02 15:04:05"
)

type LocalTime time.Time

func (date *LocalTime) Scan(value interface{}) (err error) {
	if value, ok := value.(time.Time); ok {
		*date = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", date)
}

func (date LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(date)
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

func (date LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(date)
	return []byte(fmt.Sprintf(`"%s"`, tTime.Format(LocalDateTimeFormat))), nil
}

// 实现 UnmarshalJSON 方法，将 JSON 字符串反序列化到 LocalTime 类型
func (lt *LocalTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	parsedTime, err := time.Parse(`"`+LocalDateTimeFormat+`"`, string(data))
	if err != nil {
		return err
	}
	*lt = LocalTime(parsedTime)
	return nil
}

func (lt LocalTime) String() string {
	return time.Time(lt).Format(LocalDateTimeFormat)
}
