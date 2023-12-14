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

// GormDataType gorm common data type
func (date LocalTime) GormDataType() string {
	return "time"
}

func (date *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*date)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

func (date *LocalTime) UnmarshalJSON(b []byte) error {
	return (*time.Time)(date).UnmarshalJSON(b)
}
