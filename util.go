package service

import (
	"errors"
	"github.com/satori/go.uuid"
	"math/rand"
	"strings"
	"time"
)

const (
	TimeLocation = "Asia/Shanghai"
	//TimeLayoutString        = "2006-01-02 15:04:05"
	//TimeDateLayoutString    = "2006-01-02" //only date layout
	//TimeDateIntLayoutString = "20060102"   //only date layout

)

func checkDateLocation(date1, date2 time.Time) bool {
	s := date1.Location().String()
	s2 := date2.Location().String()
	//return strings.EqualFold(date1.Location().String(), date2.Location().String())
	return strings.EqualFold(s, s2)

}

func GetDefaultLocation() *time.Location {
	location, _ := time.LoadLocation(TimeLocation)
	return location
}

//GetNow return now with default location
func GetNowWithDefaultLocation() time.Time {
	loc := GetDefaultLocation()
	return time.Now().In(loc)
}

//SubSecond returns end 減去 begin 是多少秒 ，有可能是負的
func SubSecond(begin, end time.Time) (int64, error) {

	if !checkDateLocation(begin, end) {
		return -1, errors.New("time zone not match")
	}

	return int64(end.Sub(begin).Seconds()), nil
}

func randomSerial() int {
	return randomInt(100000000000, 999999999999)
}

func randomInt(min, max int) int {
	if min >= max {
		return min
	}

	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min+1)
}

// GenerateAES256RandomBites returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateAES256RandomBytes() ([]byte, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

func uuidGetV4() string {

	u := uuid.NewV4()
	return u.String()
}
