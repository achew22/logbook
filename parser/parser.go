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
