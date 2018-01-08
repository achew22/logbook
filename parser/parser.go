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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/achew22/logbook/config"
)

var relativeDateMap = map[string]time.Duration{
	"tomorrow": 24 * time.Hour,
}

type LogEntry struct {
	Path string
	Date Date

	FutureReferences map[Date][]string
}

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

type Parser struct {
	config *config.Config

	fileMap map[Date]*LogEntry
}

func New(config *config.Config) *Parser {
	return &Parser{
		config: config,
	}
}

func (p *Parser) Parse() map[Date]*LogEntry {
	p.fileMap = map[Date]*LogEntry{}
	filepath.Walk(p.config.LogPath, p.parseFile)
	return p.fileMap
}

func YmdToDate(d string) (Date, error) {
	// Parse Year-Month-Date.
	p, err := time.Parse("2006-01-02", d)
	if err != nil {
		fmt.Printf("Unable to date (%s). err: %s\n", d, err)
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

func extractFutureReferences(d Date, b string) (map[Date][]string, error) {
	r := map[Date][]string{}

	r[TimeToDate(time.Now().Add(24*time.Hour))] = []string{
		"More text",
	}

	r[TimeToDate(time.Now())] = []string{
		"Tomorrow text",
	}

	return r, nil
}

func (p *Parser) parseFile(path string, info os.FileInfo, err error) error {
	if !strings.HasSuffix(path, ".md") {
		return nil
	}

	_, file := filepath.Split(path)
	d, err := YmdToDate(file[:len(file)-3])
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	futureReferences, err := extractFutureReferences(d, string(b))
	if err != nil {
		return err
	}

	l := &LogEntry{
		Date:             d,
		Path:             path,
		FutureReferences: futureReferences,
	}

	p.fileMap[d] = l
	return nil
}
