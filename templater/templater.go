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
package templater

import (
	"fmt"
	"os"
	"time"

	"github.com/achew22/logbook/config"
	"github.com/achew22/logbook/parser"
)

var _ = time.Now

const dateFormat = "2006-01-02"

type extractedEntry struct {
	originDate parser.Date
	text       string
}

func Print(c *config.Config, entries map[parser.Date]*parser.LogEntry, today parser.Date) string {
	out := ""
	forToday := []extractedEntry{}
	for _, entry := range entries {
		for date, lines := range entry.FutureReferences {
			if today.Equals(date) {
				for _, line := range lines {
					forToday = append(forToday, extractedEntry{
						date,
						line,
					})
				}
			}
		}
	}

	out += fmt.Sprintf("# %s - %s\n\n", c.Name, today.ToTime().Format(dateFormat))

	if len(forToday) > 0 {
		out += "Reminders:\n"

		for _, e := range forToday {
			out += fmt.Sprintf("From %s: %s\n", e.originDate.ToTime().Format(dateFormat), e.text)
		}
	}

	out += fmt.Sprintf("\n\n")

	fmt.Fprintf(os.Stderr, "Wrote file\n")
	return out
}
