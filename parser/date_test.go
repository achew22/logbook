package parser

import (
	"testing"
	"time"
)

func mustYmdToDate(d string) Date {
	s, err := YmdToDate(d)
	if err != nil {
		panic(err)
	}
	return s
}

func TestInvalidDate(t *testing.T) {
	cases := map[string]string{
		"Too many days in september": "2016-09-32",
		"Non-leapyear leap day":      "2015-02-29",
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			_, err := YmdToDate(c)
			if err == nil {
				t.Errorf("Succeeded in converting a date that should have failed: %s", c)
			}
		})
	}
}

func TestInequalDate(t *testing.T) {
	dates := map[string]struct {
		a Date
		b Date
	}{
		"YmdVsTime at midnight": {mustYmdToDate("2016-02-29"), mustYmdToDate("2016-03-01")},
	}
	for n, d := range dates {
		t.Run(n, func(t *testing.T) {
			if d.a.Equals(d.b) {
				t.Errorf("Expected a to != b in testcase. a: %v, b: %v", d.a, d.b)
			}
		})
	}
}

func TestEqualDate(t *testing.T) {
	dates := map[string]struct {
		a Date
		b Date
	}{
		"RoundTripThroughTime":    {mustYmdToDate("2011-04-02"), TimeToDate(mustYmdToDate("2011-04-02").ToTime())},
		"YmdVsManual":             {mustYmdToDate("2011-04-02"), Date{Year: 2011, Month: time.April, Day: 2}},
		"TimeVsManual":            {TimeToDate(time.Date(2012, time.July, 7, 0, 0, 0, 0, time.UTC)), Date{Year: 2012, Month: time.July, Day: 7}},
		"YmdVsTime midnight":      {mustYmdToDate("2013-03-02"), TimeToDate(time.Date(2013, time.March, 2, 0, 0, 0, 0, time.UTC))},
		"YmdVsTime midday":        {mustYmdToDate("2014-09-02"), TimeToDate(time.Date(2014, time.September, 2, 23, 59, 59, 59, time.UTC))},
		"Ymd multiple invocation": {mustYmdToDate("2015-09-02"), mustYmdToDate("2015-09-02")},
		"Time multiple invocation": {
			TimeToDate(time.Date(2014, time.September, 2, 23, 59, 59, 59, time.UTC)),
			TimeToDate(time.Date(2014, time.September, 2, 23, 59, 59, 59, time.UTC)),
		},
		"Leap day": {mustYmdToDate("2016-02-29"), Date{Year: 2016, Month: time.February, Day: 29}},
	}
	for n, d := range dates {
		t.Run(n, func(t *testing.T) {
			if !d.a.Equals(d.b) {
				t.Errorf("Expected a to = b in testcase. a: %v, b: %v", d.a, d.b)
			}
		})
	}
}
