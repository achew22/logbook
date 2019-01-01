/**
 * Copyright 2018 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the license.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specified language governing permissions and
 * limitations under the License.
 */
package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func (d Date) Equals(o Date) bool {
	return d.Year == o.Year && d.Month == o.Month && d.Day == o.Day
}

func (d Date) ToTime() time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
}

func (d Date) ToYmd() string {
	return fmt.Sprintf("%d-%02d-%02d", d.Year, d.Month, d.Day)
}

func YmdToDate(d string) (Date, error) {
	// Parse Year-Month-Date.
	p, err := time.Parse("2006-01-02", d)
	if err != nil {
		return Date{}, err
	}
	return TimeToDate(p), nil
}

func TimeToDate(t time.Time) Date {
	y, m, d := t.Date()
	return Date{
		Year:  y,
		Month: m,
		Day:   d,
	}
}

var timespecMatchers = map[*regexp.Regexp]func(Date, []string) (Date, error){
	// in X days
	// in X weeks
	// in X months
	// in X years
	regexp.MustCompile("(in )?(\\d+) (day|week|month|year)s?"): func(d Date, matches []string) (Date, error) {
		count, err := strconv.Atoi(matches[2])
		if err != nil {
			return d, fmt.Errorf("Unable to convert %q to int: %v", matches[2], err)
		}
		interval := matches[3]

		years, months, days := 0, 0, 0
		switch interval {
		case "day":
			days = count
		case "week":
			days = count * 7
		case "month":
			months = count
		case "years":
			years = count
		}
		return TimeToDate(d.ToTime().AddDate(years, months, days)), nil
	},

	// tomorrow
	regexp.MustCompile("tomorrow"): func(d Date, matches []string) (Date, error) {
		return TimeToDate(d.ToTime().AddDate(0, 0, 1)), nil
	},
}

func ParseTimespec(d Date, spec string) (Date, error) {
	spec = strings.ToLower(spec)

	for r, f := range timespecMatchers {
		if matches := r.FindStringSubmatch(spec); matches != nil {
			return f(d, matches)
		}
	}

	return d, fmt.Errorf("no valid spec parser found for spec %q", spec)
}
