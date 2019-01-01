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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	blackfriday "gopkg.in/russross/blackfriday.v2"

	"github.com/achew22/logbook/config"
)

var (
	expressionFinder = regexp.MustCompile("(.+):(.+)")
)

type ParseError struct {
	Message string `json:"message"`
}

type LogEntry struct {
	Path string
	Date Date

	PastReferences map[Date][]string

	Errors []*ParseError
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
		Errors         []*ParseError       `json:"errors,omitempty"`
	}{
		Path:           l.Path,
		Date:           l.Date.ToYmd(),
		PastReferences: marshalPastReferences(l.PastReferences),
		Errors:         l.Errors,
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
			Path:           filepath.Join(p.config.LogPath, d.ToYmd()) + ".md",
			PastReferences: map[Date][]string{},
			Errors:         []*ParseError{},
		}
	}

	return p.fileMap[d]
}

func (p *Parser) emitError(d Date, err error) {
	toLog := p.getOrCreateLog(d)
	toLog.Errors = append(toLog.Errors, &ParseError{
		Message: err.Error(),
	})
}

func (p *Parser) emitEvent(from, to Date, message string) {
	toLog := p.getOrCreateLog(to)
	toLog.PastReferences[from] = append(toLog.PastReferences[from], message)
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

	markdown := blackfriday.New()
	rootNode := markdown.Parse(b)
	rootNode.Walk(p.walkNodes(d))

	return nil
}

func (p *Parser) walkNodes(d Date) func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	return func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		switch n.Type {
		case blackfriday.Document, blackfriday.Heading, blackfriday.Paragraph:
			// These nodes can never contain any prospective information, but nodes
			// inside of them can contain info. GoToNext recurses into those nodes.
			return blackfriday.GoToNext
		case blackfriday.Text:
			p.parseEventText(d, string(n.Literal))
			return blackfriday.GoToNext
		default:
			p.emitEvent(d, d, fmt.Sprintf("Unknown node: %s %q", n, n.Literal))
			return blackfriday.GoToNext
		}
	}
}

func (p *Parser) parseEventText(d Date, text string) {
	r := expressionFinder.FindStringSubmatch(text)
	if len(r) < 3 {
		return
	}

	timespec, remark := r[1], r[2]
	reminderDate, err := ParseTimespec(d, timespec)
	if err != nil {
		p.emitError(d, err)
	}

	p.emitEvent(d, reminderDate, trim(remark))
}
