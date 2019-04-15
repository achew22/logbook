package parser

import (
	"fmt"
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

func TestParseTimespec(t *testing.T) {
	tests := map[string]struct {
		now Date
		in  string
		d   Date
		err error
	}{
		"Explicit date 2001-02-03": {mustYmdToDate("2001-02-03"), "2001-02-03", mustYmdToDate("2001-02-03"), nil},
		"Explicit date 2001-2-03":  {mustYmdToDate("2001-02-03"), "2001-2-03", mustYmdToDate("2001-02-03"), nil},
		"Explicit date 2001-02-3":  {mustYmdToDate("2001-02-03"), "2001-02-3", mustYmdToDate("2001-02-03"), nil},
		"Explicit date 2001-2-3":   {mustYmdToDate("2001-02-03"), "2001-2-3", mustYmdToDate("2001-02-03"), nil},
		"Tomorrow":                 {mustYmdToDate("2001-02-03"), "tomorrow", mustYmdToDate("2001-02-04"), nil},
		"case sensitive tomorrow":  {mustYmdToDate("2001-02-03"), "TOMorrow", mustYmdToDate("2001-02-04"), nil},
		"Future date":              {mustYmdToDate("2019-03-11"), "2019-05-14", mustYmdToDate("2019-05-14"), nil},

		"Error input": {mustYmdToDate("2001-02-03"), "82872--1", mustYmdToDate("2001-02-03"), fmt.Errorf("no valid spec parser found for spec \"82872--1\"")},
	}
	for n, test := range tests {
		t.Run(n, func(t *testing.T) {
			r, err := ParseTimespec(test.now, test.in)
			if err == nil && test.err != nil {
				t.Fatalf("Expected an error but didn't get one")
			}
			if err != nil && test.err == nil {
				t.Fatalf("Expected no error but got [%v]", err)
			}
			if err != nil && err.Error() != test.err.Error() {
				t.Fatalf("Expected error to be [%v], was [%v]", test.err, err)
			}
			if !test.d.Equals(r) {
				t.Fatalf("Expected date to be [%v], was [%v]", r, test.d)
			}
		})
	}
}
