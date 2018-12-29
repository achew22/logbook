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
	"encoding/json"
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

	PastReferences map[Date][]string
}

func marshalPastReferences(r map[Date][]string) map[string][]string {
	out := map[string][]string{}
	for k, v := range r {
		out[k.ToYmd()] = v
	}
	return out
}
func (l *LogEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Path           string              `json:"path"`
		Date           string              `json:"date"`
		PastReferences map[string][]string `json:"pastReferences"`
	}{
		Path:           l.Path,
		Date:           l.Date.ToYmd(),
		PastReferences: marshalPastReferences(l.PastReferences),
	})
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

	// First walk all the files in thie directory extracting any forward looking information they might have.
	filepath.Walk(p.config.LogPath, p.parseFile)

	return p.fileMap
}

// getOrCreateLog either gets from the eisting fileMap a date or creates it.
func (p *Parser) getOrCreateLog(d Date) *LogEntry {
	_, ok := p.fileMap[d]
	if !ok {
		p.fileMap[d] = &LogEntry{
			Date:           d,
			Path:           filepath.Join(p.config.LogPath, d.ToYmd()),
			PastReferences: map[Date][]string{},
		}
	}

	return p.fileMap[d]
}

func (p *Parser) emitEvent(from, to Date, message string) error {
	toLog := p.getOrCreateLog(to)
	toLog.PastReferences[from] = append(toLog.PastReferences[from], message)

	return nil
}

// emitFutureReferences parses the string of the content for important
// dates and returns a list of dates/events that this entry created.
func (p *Parser) emitFutureReferences(d Date, content string) error {
	err := p.emitEvent(d, TimeToDate(d.ToTime().AddDate(0, 0, 5)), "More text")
	if err != nil {
		return err
	}

	// Pretend we pared out "Tomorrow: More text"
	err = p.emitEvent(d, TimeToDate(d.ToTime().AddDate(0, 0, 1)), "Tomorrow text")
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseFile(path string, info os.FileInfo, err error) error {
	if !strings.HasSuffix(path, ".md") {
		return nil
	}

	name := filepath.Base(path)
	d, err := YmdToDate(name[:len(name)-len(".md")])
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return p.emitFutureReferences(d, string(b))
}
