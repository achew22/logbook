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
	return time.Date(d.Year, time.Month(d.Month), d.Year, 0, 0, 0, 0, time.UTC)
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
